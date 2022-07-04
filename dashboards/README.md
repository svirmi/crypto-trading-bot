## Interactive visualization with Streamlit

We suggest installing [Miniconda](https://docs.conda.io/en/latest/miniconda.html) or [Miniforge](https://github.com/conda-forge/miniforge).

Alternative vitrual environment libraries: [pipenv](https://pypi.org/project/pipenv/), [venv](https://docs.python.org/3/library/venv.html)

### 1. Setup virtual python environment

Create a virtual python environment with miniconda (you can use `venv` instead)

```
$ conda create -n crypto-bot python=3.9 -y
```
Activate the vistual environment

```
$ conda activate crypto-bot
```

Install required packages

```
$ pip install -r requirements.txt
```

### 2. Run the streamlit app

To visualize the interactive dashboards on you browser run

```
$ streamlit run dashboards.py 
```


### 3. TODO Explain the dashboards