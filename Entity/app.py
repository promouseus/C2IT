"""C2 Data Analysis"""
import numpy as np
import pandas as pd

#FastAPI for future:
# https://learn.vonage.com/blog/2021/08/10/the-ultimate-face-off-flask-vs-fastapi/

# s = pd.Series([1, 3, 5, np.nan, 6, 8])

# print(s)

# https://github.com/RichDijk/EAGOC/blob/master/ea_goc.ipynb

modelElements = pd.read_xml('ModelExchange/case1.xml', xpath="//archimate:elements/*", namespaces={"archimate": "http://www.opengroup.org/xsd/archimate/3.0/"})
modelRelationships = pd.read_xml('ModelExchange/case1.xml', xpath="//archimate:relationships/*", namespaces={"archimate": "http://www.opengroup.org/xsd/archimate/3.0/"})
modelOrganizations = pd.read_xml('ModelExchange/case1.xml', xpath="/archimate:model/archimate:organizations/*", namespaces={"archimate": "http://www.opengroup.org/xsd/archimate/3.0/"})

modelProperty = pd.read_xml('ModelExchange/case1.xml', xpath="/archimate:model/archimate:propertyDefinitions/*", namespaces={"archimate": "http://www.opengroup.org/xsd/archimate/3.0/"})

#https://www.opengroup.org/xsd/archimate/3.1/html-model/archimate3_Model.html#PropertyDefinitionType Now string but more types avaiable

print("This line will be printed.")

# https://code.visualstudio.com/docs/python/python-tutorial
# python3 -m venv .venv (from Workspace root)
# source .venv/bin/activate
