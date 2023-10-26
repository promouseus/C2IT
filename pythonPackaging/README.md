# C2IT agent
This is the C2IT agent written in Python


## Development notes:

```bash
https://packaging.python.org/en/latest/tutorials/packaging-projects/
source .venv/bin/activate
python3 -m build
python3 -m twine upload --repository testpypi dist/*

docker run -it --entrypoint bash python:3.5.10-slim
```

## Run from source dir
```bash
src/
python -m c2it.agent
```
## Install from test PyPi
```bash
pip install -i https://test.pypi.org/simple/ c2it
alias c2it="python -m c2it"
```
