from sqlalchemy import PrimaryKeyConstraint
from app import db
from datetime import datetime

class Test(db.Model):
    __tablename__ = 'test'

    id = db.Column(db.Integer,primary_key=True)

















