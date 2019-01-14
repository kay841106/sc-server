
import sys
sys.path.append('/home/babyuser/Documents/project/smartUniversity')

import pandas as pd


from pymongo import MongoClient
from pymongo import DESCENDING as descend
from pymongo import ASCENDING as ascend
from datetime import datetime, timedelta, date
import collections
import logging
# from operator import itemgette
import json
import pymongo

def dbconnect():
    return pymongo.MongoClient('mongodb://140.118.70.136:10003/')

_PATH = './DOC/devices.csv'
_DB_COLL = 'SC01_DeviceManager'
#READ DATA FROM CSV
def input_data(path):
    data1 = pd.read_csv(_PATH, header=1)
    asu = data1.drop(data1.index[:0])
    asu3 = asu.loc[:, ~asu.columns.str.contains('^Unnamed')]
    asu2 = asu3.fillna(method='ffill')
    return(asu2)


def connect_db():
    conn = dbconnect()
    db = conn['Bimo_test']
    db.authenticate('dontask', 'idontknow','admin')
    srcB = db[_DB_COLL]
    return(srcB)

def proc_insert():
    tmp_listID = []
    dstA = connect_db()
    the_list = input_data(_PATH)
    print(len(the_list))

    for cnt in range(len(the_list)):
        # for the_list in the_list:
        # print(the_list)
        tmpdevID1 = the_list.MAC[cnt]
        tmpdevID2 = the_list.ID2[cnt]  
        devID = '330005'+tmpdevID1+str(tmpdevID2)
        # print(type(tmpdevID2))
        _GWID = the_list.GWID[cnt][0:14]                                                                                                                                                                                                                            
        record = {
            'Device_Brand':'AAEON',
            'Device_Type': the_list.TYPE[cnt],
            'Building_Details':the_list.PLACE[cnt],
            'Building_Name':the_list.Building_Name[cnt],
            'devID':devID,
            'GWID':_GWID,
            'Device_Name':the_list.NUM[cnt],
            'Device_Details':the_list.TERRITORY[cnt],
            'Floor':the_list.FLOOR[cnt],
            'Device_Info':the_list.METER_EN[cnt]
            
        }
        # print("UPDATE {}".format(cnt))
        
        tmp_listID.append(devID)
        dstA.update_one({'devID':devID},{'$set':record},upsert=True)
        # dstA.insert_one(record).acknowledged
    return tmp_listID    

def counter_feiting(the_list):
    # counter_list = 
    # print(the_list)
    print([item for item, count in collections.Counter(the_list).items() if count > 1])

def main():
    letwork = input_data(_PATH)
    print(letwork)
    connect_db()
    counter_feiting(proc_insert())

if(__name__ == "__main__"):
    main()
