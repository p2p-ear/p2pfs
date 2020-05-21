# project/server/models.py

import datetime
import jwt # for encode_auth_token method in class User
import json

from project.server import app, db, bcrypt

class User(db.Model):
    """ User Model for storing user related details """
    __tablename__ = "users"

    id = db.Column(db.Integer, primary_key=True, autoincrement=True)
    email = db.Column(db.String(255), unique=True, nullable=False)
    password = db.Column(db.String(255), nullable=False)
    default_data = {
                "Name": "/", 
                "IsDir": True, 
                "Size": 0, 
                "Flag":  False,
                "Child": 
                {
                    "Uploaded": {
                        "Name": "Uploaded", 
                        "IsDir": True, 
                        "Size": 0, 
                        "Flag":  False, 
                        "Child": 
                        {
                        }
                    }
                }
            }
    data = db.Column(db.Text, nullable=False, default=json.dumps(default_data)) # For JSON
    files = db.relationship('File', backref='owner')
    coins = db.Column(db.Integer, nullable=False, default=0) # For storing money
    registered_on = db.Column(db.DateTime, nullable=False, default=datetime.datetime.utcnow()) # Added: default=datetime.datetime.utcnow()
    admin = db.Column(db.Boolean, nullable=False, default=False)

    def __init__(self, email, password, admin=False):
        self.email = email
        self.password = bcrypt.generate_password_hash(
            password, app.config.get('BCRYPT_LOG_ROUNDS')
        ).decode()

    
    def encode_auth_token(self, user_id):
        """
        Generates the Auth Token
        :return: string
        """
        try:
            payload = {
                'exp': datetime.datetime.utcnow() + datetime.timedelta(days=0, hours=1),
                'iat': datetime.datetime.utcnow(),
                'sub': user_id
            }
            return jwt.encode(
                payload,
                app.config.get('SECRET_KEY'),
                algorithm='HS256'
            )
        except Exception as e:
            return e

        
    @staticmethod
    def decode_auth_token(auth_token):
        """
        Validates the auth token
        :param auth_token:
        :return: integer|string
        """
        try:
            payload = jwt.decode(auth_token, app.config.get('SECRET_KEY'))
            is_blacklisted_token = BlacklistToken.check_blacklist(auth_token)
            if is_blacklisted_token:
                return 'Token blacklisted. Please log in again.'
            else:
                return payload['sub']
        except jwt.ExpiredSignatureError:
            return 'Signature expired. Please log in again.'
        except jwt.InvalidTokenError:
            return 'Invalid token. Please log in again.'


class File(db.Model):
    """ File Model for storing user related files """
    __tablename__ = "files"

    id = db.Column(db.Integer, primary_key=True)
    user_id = db.Column(db.Integer, db.ForeignKey('users.id'), nullable=False)
    file_name = db.Column(db.Text, nullable=False)
    file_size = db.Column(db.BigInteger, nullable=False)
    total_shards = db.Column(db.BigInteger, nullable=False)
    initial_ip = db.Column(db.BigInteger, nullable=False)

    def __init__(self, user_id, file_name, file_size, n_shards, initial_ip):
        self.user_id = user_id
        self.file_name = file_name
        self.file_size = file_size
        self.total_shards = n_shards
        self.initial_ip = initial_ip


    def __repr__(self):
        return "File('{}', '{}')".format(self.user_id, self.file_name)


class Node(db.Model):
    """ Node Model for storing node related IPs """
    __tablename__ = "nodes"

    id = db.Column(db.Integer, primary_key=True)
    ip_address = db.Column(db.BigInteger, nullable=False)

    def __init__(self, ip_address):
        self.ip_address = ip_address
    
    def __repr__(self):
        return "File('{}')".format(self.ip_address)


class BlacklistToken(db.Model):
    """
    Token Model for storing JWT tokens
    """
    __tablename__ = 'blacklist_tokens'

    id = db.Column(db.Integer, primary_key=True, autoincrement=True)
    token = db.Column(db.String(500), unique=True, nullable=False)
    blacklisted_on = db.Column(db.DateTime, nullable=False)

    def __init__(self, token):
        self.token = token
        self.blacklisted_on = datetime.datetime.now()

    def __repr__(self):
        return '<id: token: {}'.format(self.token)

    @staticmethod
    def check_blacklist(auth_token):
        # check whether auth token has been blacklisted
        res = BlacklistToken.query.filter_by(token=str(auth_token)).first()
        if res:
            return True  
        else:
            return False


class Сertificate(db.Model):
    """
    Сertificate Model for storing JWT tokens
    """
    __tablename__ = 'certificates'

    id = db.Column(db.Integer, primary_key=True, autoincrement=True)
    token = db.Column(db.String(500), unique=True, nullable=False)
    shards = db.Column(db.Integer, nullable=False)

    def __init__(self, user_id, mini_shard_size, shards, act, file_name):
        self.token = self.encode_certificate_token(user_id, mini_shard_size, act, file_name).decode()
        self.shards = shards


    def __repr__(self):
        return '<id: token: {}'.format(self.token)

    @staticmethod
    def encode_certificate_token(user_id, mini_shard_size, act, file_name):
        """
        Generates the certificate Token
        :return: string
        """
        try:
            payload = {
                'exp': datetime.datetime.utcnow() + datetime.timedelta(days=1, seconds=0),
                'iat': datetime.datetime.utcnow(),
                'sub': user_id,
                'size': mini_shard_size,
                'act': act,
                'name': file_name
            }
            return jwt.encode(
                payload,
                app.config.get('SECRET_KEY'),
                algorithm='HS256'
            )
        except Exception as e:
            return e

    @staticmethod
    def _decode_certificate_token(certificate_token):
        """
        Validates the certificate token
        :param certificate_token:
        :return: integer|string
        """
        try:
            payload = jwt.decode(auth_token, app.config.get('SECRET_KEY'))
            is_blacklisted_token = BlacklistToken.check_blacklist(certificate_token)
            if is_blacklisted_token:
                return 'Token blacklisted. Please request again.'
            else:
                return payload['act']#payload['sub']
        except jwt.ExpiredSignatureError:
            return 'Signature expired. Please request again.'
        except jwt.InvalidTokenError:
            return 'Invalid token. Please request again.'