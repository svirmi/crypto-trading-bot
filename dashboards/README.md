# Simulation visualization with Streamlit
Dashboard created for visualizing data about a simulation in order to evaluate a strategy's performance, be able to compare strategies with one another and extract insight for developing new strategies. 

## Requiremets
In order to run the dashboard, a tool to virtualize a python environment is required. Any of the tools listed below will be fine but Miniconda is the suggested one.
- [miniconda](https://docs.conda.io/en/latest/miniconda.html) (suggested environment)
- [miniforge](https://github.com/conda-forge/miniforge)
- [pipenv](https://pypi.org/project/pipenv/)
- [venv](https://docs.python.org/3/library/venv.html)

## Setup virtual python environment
Create a virtual python environment with you tool of choise. Instructions are provided here for Miniconda only:
```
$ conda create -n crypto-trading-bot python=3.9 -y
```
Activate the virtual environment:
```
$ conda activate crypto-trading-bot
```
Install required packages
```
$ pip install -r requirements.txt
```

## Run the streamlit app
To visualize the interactive dashboards on you browser run
```
$ streamlit run dashboards.py 
```

## TODO: dashboard documentation