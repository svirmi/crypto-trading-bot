from audioop import avg
from distutils import errors
import streamlit as st
from pymongo import MongoClient
import pandas as pd
from datetime import datetime
import plotly.express as px
import math
from plotly import graph_objects as go

MONGO_HOST = 'localhost'  # TODO
MONGO_PORT = 27017  # TODO

MONGO_DB = 'ctb-simulation'
MONGO_COLLECTION = 'analytics' 
EXECUTION = 'execution-analytics'
OPERATIONS = 'operation-analytics'
WALLET = 'wallet-analytics'

st.set_page_config(layout="wide")

st.title('Simulation Dashboards')
st.write('-'*20) 

# -----------------------------------------------------------------------------------------------------------------
# LOAD DATA 
# -----------------------------------------------------------------------------------------------------------------

# Connect to mongodb collection analytics
mongo_client = MongoClient(MONGO_HOST, MONGO_PORT)
analytics = mongo_client[MONGO_DB][MONGO_COLLECTION]

# Load finished executions 
executions = list(analytics.find({'analyticsType': EXECUTION, 'status': 'EXE_TERMINATED'}))

strategy_types = []
for execution in executions:
    if (strategy := execution['strategyType']) not in strategy_types:
        strategy_types.append(strategy)

# Create selection box with trading strategies
strategy = st.selectbox('<Select a trading strategy>', strategy_types)
# TODO: add strategy description

# Create selection box with execution IDs
exe_id_assets = st.selectbox('<Select an execution ID>', 
                      (f'{e["exeId"]} {"@".join(e["assets"])}' for e in executions if e['strategyType'] == strategy))
exe_id = exe_id_assets.split(' ')[0]
# full assets
assets = exe_id_assets.split(' ')[1].split('@')
# crypto assets
crypto_assets = [a for a in assets if a != 'USDT']

# Load wallet and operations
wallet = list(analytics.find({'exeId': exe_id, 'analyticsType': WALLET}))
wallet = pd.json_normalize(wallet)
wallet = wallet.drop(columns=['_id', 'exeId', 'analyticsType'] + [f'assetStatuses.{ast}.asset' for ast in assets])
wallet = wallet.rename(columns={c: '_'.join(c.split('.')[1:]) for c in wallet.columns if 'assetStatuses' in c})
wallet['timestamp'] = wallet['timestamp'].apply(lambda d: datetime.fromtimestamp(d / 1000000))
wallet = wallet.sort_values('timestamp').reset_index(drop=True)
# convert to float
for col in wallet.columns:
    if col != 'timestamp':
        wallet[col] = wallet[col].astype(float)

operations = list(analytics.find({'exeId': exe_id, 'analyticsType': OPERATIONS}))
operations = pd.json_normalize(operations)
operations = operations.drop(columns=['_id', 'exeId', 'analyticsType'])
operations['timestamp'] = operations['timestamp'].apply(lambda d: datetime.fromtimestamp(d / 1000000))
operations = operations.sort_values('timestamp').reset_index(drop=True)
for col in ['amount', 'price']:
    operations[col] = operations[col].astype(float)
operations.loc[operations['amountSide'] == 'QUOTE_AMOUNT', 'amount'] = \
    operations.loc[operations['amountSide'] == 'QUOTE_AMOUNT', 'amount'] / \
    operations.loc[operations['amountSide'] == 'QUOTE_AMOUNT', 'price']
operations.loc[operations['side'] == 'SELL', 'amount'] *= -1
operations = operations.drop(columns=['amountSide'])

# -----------------------------------------------------------------------------------------------------------------
# SIMULATION PARAMETERS
# -----------------------------------------------------------------------------------------------------------------

# Show strategy parameters, initial and final values
col1, col2, col3 = st.columns(3)
# 1. strategy params
selected_execution = next((e for e in executions if e['exeId'] == exe_id))
strategy_params = '\n\n'.join([key + ':' + (25-len(key))*' ' + val 
                             for key, val in selected_execution['props'].items()])
with col1:
    st.subheader('Strategy parameters')
    st.text(strategy_params)

def get_amount_and_prices(index):
    amounts = {asset: wallet[f'{asset}_amount'][index] for asset in assets}
    prices = {asset: wallet[f'{asset}_price'][index] for asset in assets}
    value_string = ''
    for asset in assets:
        value_string += \
            f'{asset}:\n\tamount:' + 10*' ' + str(amounts[asset]) + \
            f'\n\tprice:' + 11*' ' + str(prices[asset]) + '\n'
    value_string += '\nTotal value ($):' + 9*' ' + str(math.floor(sum([prices[a]*amounts[a] for a in assets])))
    return value_string

with col2:
    st.subheader('Initial amounts and prices')
    st.text(get_amount_and_prices(0))

with col3:
    st.subheader('Final amounts and prices')
    st.text(get_amount_and_prices(len(wallet)-1))


st.write('-'*20)  
# -----------------------------------------------------------------------------------------------------------------
# WALLET VS BASELINE
# -----------------------------------------------------------------------------------------------------------------

st.subheader('Wallet value using strategy vs baseline')

initial_USDT = wallet['USDT_amount'][0]
wallet['baseline'] = initial_USDT
for asset in crypto_assets:
    wallet['baseline'] += wallet[f'{asset}_amount'][0] * wallet[f'{asset}_price']
wallet['relativeGain(%)'] = (wallet['walletValue'] - wallet['baseline'])/ wallet['baseline'] * 100

col1, col2 = st.columns(2)
with col1:
    fig = px.line(wallet, x='timestamp', y=['walletValue', 'baseline'], 
                  title='Wallet value ($) vs baseline (invest all at time 0)')
    fig.update_layout(legend=dict(yanchor="top", y=0.99, xanchor="left", x=0.01))
    fig.add_hline(y=wallet['baseline'][0], line_width=3, line_dash="dash", line_color="green")
    st.plotly_chart(fig)
with col2:
    fig = px.line(wallet, x='timestamp', y='relativeGain(%)',
                  title='Relative gain in percentage with respect to baseline')
    fig.add_hline(y=0, line_width=3, line_dash="dash", line_color="red")
    fig.add_hrect(y0=0, y1=wallet['relativeGain(%)'].min(), 
                  line_width=0, fillcolor="red", opacity=0.2)
    fig.add_hrect(y0=0, y1=wallet['relativeGain(%)'].max(), 
                  line_width=0, fillcolor="green", opacity=0.2)
    st.plotly_chart(fig)

st.write('-'*20)  
# -----------------------------------------------------------------------------------------------------------------
# OPERATIONS
# -----------------------------------------------------------------------------------------------------------------

st.subheader('Crypto prices and operations.')

def get_operation_segment(operation, low_val):
        return go.layout.Shape(type="line", x0=operation['timestamp'], x1=operation['timestamp'],
                            y0=low_val, y1=operation['price'], 
                            line=dict(color='red' if operation['side'] == 'SELL' else 'green', width=2, dash='dot'), 
                            opacity=0.8, layer="below")  # , annotation=operation['amount'])

for asset in crypto_assets:
    
    buys = operations[(operations['side']=='BUY') & (operations["base"] == asset)]
    sells = operations[(operations['side']=='SELL') & (operations["base"] == asset)]

    nsells = len(sells)
    nbuys = len(buys)
    avg_buy_price = (buys['amount'] * buys['price']).sum() /  buys['amount'].sum()
    avg_sell_price = (sells['amount'] * sells['price']).sum() /  sells['amount'].sum()
    operations_summary = \
        f"""
        #### {asset}  \n  \n

        Number of **SELL** operations: {nsells}  \n 
        Number of **BUY** operations: {nbuys}  \n
        """
    avg_price_summary = \
        f"""
        ***  \n  \n

        Average buy price: {"%.2f" % avg_buy_price} $ \n
        Average sell price: {"%.2f" % avg_sell_price} $ \n
        """

    col1, col2 = st.columns(2)
    with col1:
        st.markdown(operations_summary)

    with col2:
        st.markdown(avg_price_summary)

    fig = px.line(wallet, x='timestamp', y=f'{asset}_price')
    ops = operations[operations['base'] == asset]
    ops['amount'] = ops['amount'].apply(lambda a: "%.2f" % float(a))
    color_discrete_map = {'SELL': 'rgb(255,0,0)', 'BUY': 'rgb(0,255,0)'}
    fig2 = px.scatter(ops, x='timestamp', y='price', color='side', 
                    color_discrete_map=color_discrete_map, size_max=2, 
                    text='amount')
    fig = go.Figure(data=fig.data+fig2.data)
    low_val = wallet[f'{asset}_price'].min()
    ply_shapes = []
    for _, op in operations[operations['base'] == asset].iterrows():
        ply_shapes.append(get_operation_segment(op, low_val))
    fig.update_layout(shapes=ply_shapes)
    fig.update_layout(legend=dict(yanchor="top", y=0.99, xanchor="left", x=0.01))
    fig.update_xaxes(title='timestamp')
    fig.update_yaxes(title=f'{asset}_price')
        
    st.plotly_chart(fig, use_container_width=True)

        
# Streamlit widgets automatically run the script from top to bottom. Since
# this button is not connected to any other logic, it just causes a plain rerun.
st.button("Re-run")
