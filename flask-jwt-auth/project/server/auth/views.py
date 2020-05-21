# project/server/auth/views.py

from flask import Blueprint, request, make_response, jsonify
from flask.views import MethodView

from project.server import bcrypt, db
from project.server.models import User, BlacklistToken, Сertificate, File, Node

import json
from functools import reduce  # forward compatibility for Python 3
import operator
from werkzeug.exceptions import BadRequest
from ipaddress import ip_address as IPA
import math

auth_blueprint = Blueprint('auth', __name__)

class RegisterAPI(MethodView):
    """
    User Registration Resource
    """

    def post(self):
        try:
            # Get the post data
            post_data = request.get_json()

            if not post_data: # Request isn't JSON type
                raise BadRequest

            # Check if data is incorrect
            if not post_data.get('email') or not isinstance(post_data.get('email'), str):
                raise ValueError({'status': 1, 'message': 'You have forgotten to specify email or its type is incorrect!'}, 400)
            if not post_data.get('password') or not isinstance(post_data.get('password'), str):
                raise ValueError({'status': 1, 'message': 'You have forgotten to specify password or its type is incorrect!'}, 400)
            #if (post_data.get('body') != dict()):
            #    raise ValueError({'status': 1, 'message': 'You have forgotten to specify body!'}, 400)
            
            # Check if user already exists
            user = User.query.filter_by(email=post_data.get('email')).first()
            if not user:
                user = User(
                    email=post_data.get('email'),
                    password=post_data.get('password')
                )

                # Insert the user
                db.session.add(user)
                db.session.commit()
                # Generate the auth token
                auth_token = user.encode_auth_token(user.id)
                responseObject = {
                    'status': 0,
                    'message': 'Successfully registered.',
                    'auth_token': auth_token.decode()
                }
                return make_response(jsonify(responseObject)), 201
            else:
                responseObject = {
                    'status': 1,
                    'message': 'User already exists. Please Log in.',
                }
                return make_response(jsonify(responseObject)), 202
        
        except ValueError as responseObject:
            return make_response(jsonify(responseObject.args[0])), responseObject.args[1]

        except BadRequest:
            return make_response(jsonify({'status': 1, 'message': 'Request should be JSON type!'})), 400

        except Exception as e:
            return make_response(jsonify({'status': 1, 'message': 'Some error occurred. Please try again.'})), 401


class LoginAPI(MethodView):
    """
    User Login Resource
    """
    def post(self):
        try:
            # Get the post data
            post_data = request.get_json()

            if not post_data: # Request isn't JSON type
                raise BadRequest

            # Check if data is incorrect
            if not post_data.get('email') or not isinstance(post_data.get('email'), str):
                raise ValueError({'status': 1, 'message': 'You have forgotten to specify email or its type is incorrect!'}, 400)
            if not post_data.get('password') or not isinstance(post_data.get('password'), str):
                raise ValueError({'status': 1, 'message': 'You have forgotten to specify password or its type is incorrect!'}, 400)
            if (post_data.get('body') != {}):
                raise ValueError({'status': 1, 'message': 'You have forgotten to specify body!'}, 400)
            if len(post_data) != 3:
                raise ValueError({'status': 1, 'message': 'Too many arguments!'}, 400)
        
            # Fetch the user data
            user = User.query.filter_by(
                email=post_data.get('email')
            ).first()
            if user and bcrypt.check_password_hash(
                user.password, post_data.get('password')
            ):
                auth_token = user.encode_auth_token(user.id)
                if auth_token:
                    responseObject = {
                        'status': 0,
                        'message': 'Successfully logged in.',
                        'email': post_data.get('email'),
                        'auth_token': auth_token.decode()
                    }
                    return make_response(jsonify(responseObject)), 200
            elif user and not bcrypt.check_password_hash( # If incorrect password
                user.password, post_data.get('password')
            ):
                return make_response(jsonify({'status': 1, 'message': 'Incorrect password.'})), 404
            else:
                return make_response(jsonify({'status': 1, 'message': 'User does not exist.'})), 404
        
        except ValueError as responseObject:
            return make_response(jsonify(responseObject.args[0])), responseObject.args[1]

        except BadRequest:
            return make_response(jsonify({'status': 1, 'message': 'Request should be JSON type!'})), 400

        except Exception as e:
            return make_response(jsonify({'status': 1, 'message': 'Try again'})), 500
            

class UserAPI(MethodView):
    """
    User Resource
    """
    def get(self):
        # get the auth token
        auth_header = request.headers.get('Authorization')
        if auth_header:
            try:
                auth_token = auth_header.split(" ")[1]
            except IndexError:
                responseObject = {
                    'status': 1,
                    'message': 'Bearer token malformed.'
                }
                return make_response(jsonify(responseObject)), 401
        else:
            auth_token = ''
        if auth_token:
            resp = User.decode_auth_token(auth_token)
            if not isinstance(resp, str):
                user = User.query.filter_by(id=resp).first()
                responseObject = {
                    'status': 0,
                    'data': {
                        'user_id': user.id,
                        'email': user.email,
                        'admin': user.admin,
                        'registered_on': user.registered_on,
                        'coins': user.coins
                    }
                }
                return make_response(jsonify(responseObject)), 200
            return make_response(jsonify({'status': 1, 'message': resp})), 401
        else:
            return make_response(jsonify({'status': 1, 'message': 'Provide a valid auth token.'})), 401


class LogoutAPI(MethodView):
    """
    Logout Resource
    """
    def post(self):
        # Get auth token
        auth_header = request.headers.get('Authorization')
        if auth_header:
            auth_token = auth_header.split(" ")[1]
        else:
            auth_token = ''
        if auth_token:
            resp = User.decode_auth_token(auth_token)
            if not isinstance(resp, str):
                # mark the token as blacklisted
                blacklist_token = BlacklistToken(token=auth_token)
                try:
                    # insert the token
                    db.session.add(blacklist_token)
                    db.session.commit()
                    return make_response(jsonify({'status': 0, 'message': 'Successfully logged out.'})), 200
                except Exception as e:
                    return make_response(jsonify({'status': 1, 'message': e})), 200
            else:
                return make_response(jsonify({'status': 1, 'message': resp})), 401
        else:
            return make_response(jsonify({'status': 1, 'message': 'Provide a valid auth token.'})), 403


class RequestAPI(MethodView):
    """
    User Update Resource
    """
    page_size = 4096
    num_pages = 1600
    shard_size = num_pages * page_size
    num_mini_shards = 8
    mini_shard_size = shard_size / num_mini_shards


    def get_body(self):
        """Get body from request JSON."""
        self.body = self.post_data.get('body')
        if not self.body: # Empty 'body'
            if not ((self.post_data.get('type') == 3) or (self.post_data.get('type') == 4)):
                raise ValueError({'status': 1, 'message': 'You should specify body!'}, 400)
        if not isinstance(self.body, type({})):
            raise ValueError({'status': 1, 'message': 'Body must be a dictionary type!'}, 400)
        if (not self.post_data.get('pass') or not self.post_data.get('email') or not self.post_data.get('JWT')):
            raise ValueError({'status': 1, 'message': 'You have forgotten to specify login, password or JWT!'}, 400)

        # Selecting user
        self.user = User.query.filter_by(email=self.post_data.get('email')).first()
        if not self.user: # User doesn't exist
            raise ValueError({'status': 1, 'message': 'User doesn\'t exist!'}, 404)


    def get_by_path(self, root, items):
        """Access a nested object in root by item sequence."""
        return reduce(operator.getitem, items, root)


    def set_item_by_path(self, root, items, value, item):
        """Set a directory or file in a nested object in root by item sequence."""
        if value in self.get_by_path(root, items):
            raise ValueError({'status': 1, 'message': 'File or directory already exists!'}, 400)
        if self.get_by_path(root, items[:-1])["Flag"]:
            raise ValueError({'status': 1, 'message': 'You can\'t add file or directory into real directory!'}, 403)
        self.get_by_path(root, items)[value] = item

    
    def delete_by_path(self, root, items, value):
        """Delete a directory or file in a nested object in root by item sequence."""
        del self.get_by_path(root, items)[value]


    def add_item(self, item, request_type):
        """Add item into directory tree."""
        if not self.body.get('name'):
            raise ValueError({'status': 1, 'message': 'You should specify file or directory name!'}, 400)
        if not self.body.get('path'):
            raise ValueError({'status': 1, 'message': 'You should specify path!'}, 400)
        
        path = self.body.get('path')    
        if (self.body.get('name').find("/") != -1):
            raise ValueError({'status': 1, 'message': 'You can\'t use "/" symbol in directory name!'}, 400)
        if (path[0] != '/'):
            raise ValueError({'status': 1, 'message': 'You must always start your path from "/" symbol!'}, 400)
        data = json.loads(self.user.data)
        if (path[-1] == '/'):
            initial_path = path.split('/')[1:-1]
        else:
            initial_path = path.split('/')[1:]
        
        initial_path.append(self.body.get('name'))
        abs_path = ["Child"]
        for i in initial_path[:-1]:
            abs_path.append(i)
            abs_path.append("Child")
        
        item["Name"] = initial_path[-1]
        responseObject = {
            'status': 0,
            'type': request_type,
            'message': 'You have successfully added new file or directory!',
            'email': self.post_data.get('email'),
            'body': {}
        }
        self.set_item_by_path(data, abs_path, initial_path[-1], item)
        self.user.data = json.dumps(data)
        if (request_type == 5):
            name = path + self.body.get('name')
            exist_file = File.query.filter_by(file_name=name).first()
            if exist_file: # File already exists
                raise ValueError({'status': 1, 'message': 'File or directory already exists!'}, 400)

            n_shards = (item["Size"] // self.shard_size) * self.num_mini_shards + math.ceil((item["Size"] % self.shard_size) / self.mini_shard_size)
            certificate = Сertificate(user_id=self.user_id, mini_shard_size=self.mini_shard_size, shards=n_shards, act=1, file_name=name)
            new_file = File(user_id=self.user_id, file_name=name, file_size=item["Size"])
                    
            if certificate and new_file:
                responseObject['body'] = {
                    'certificate_token': certificate.token
                }
                db.session.add(certificate)
                db.session.add(new_file)
                
                self.user.coins -= item["Size"]
                                
        db.session.commit()
        return make_response(jsonify(responseObject)), 200

    def post(self):
        try:
            # get the post data
            self.post_data = request.get_json()

            if not self.post_data: # Request isn't JSON type
                raise BadRequest

            # Get auth token
            auth_token = self.post_data.get('JWT')
            
            if auth_token:
                self.user_id = User.decode_auth_token(auth_token)
                if not isinstance(self.user_id, str):
                    if (self.post_data.get('type') == 0): # Add abstract directory
                        self.get_body()
                        item = {
                            "Size": 0,
                            "IsDir": True,
                            "Flag": False,
                            "Child": {}
                            }
                        return self.add_item(item, 0)           
                    elif (self.post_data.get('type') == 1): # Delete
                        self.get_body()
                        if not (self.body.get('path') and (len(self.body) == 1)):
                            raise ValueError({'status': 1, 'message': 'You should specify path in "Delete" method!'}, 400)
                        
                        path = self.body.get('path')
                        abs_path = ["Child"]
                        if (path[0] != '/'):
                            raise ValueError({'status': 1, 'message': 'You must always start your path from "/" symbol!'}, 400)
                        if (len(path) == 1):
                            raise ValueError({'status': 1, 'message': 'You must specify correct path!'}, 400)

                        data = json.loads(self.user.data)
                        if (path[-1] == '/'):
                            initial_path = path.split('/')[1:-1]
                        else:
                            initial_path = path.split('/')[1:]
                        
                        abs_path = ["Child"]
                        for i in initial_path[:-1]:
                            abs_path.append(i)
                            abs_path.append("Child")
                        
                        self.delete_by_path(data, abs_path, initial_path[-1])
                        self.user.data = json.dumps(data)
                        db.session.commit()
                        
                        responseObject = {
                            'status': 0,
                            'type': 1,
                            'message': 'You have successfuly deleted!',
                            'email': self.user.email,
                            'body': {}
                        }
                        return make_response(jsonify(responseObject)), 200
                    elif (self.post_data.get('type') == 2): # Add coins
                        self.get_body()
                        if isinstance(self.body.get('value'), int) and (len(self.body) == 1):
                            if (self.body.get('value') <= 0):
                                raise ValueError({'status': 1, 'message': "You can't add negative number of coins!"}, 400)
                            self.user.coins += self.body.get('value')
                            db.session.commit()
                            responseObject = {
                                'status': 0,
                                'type': 2,
                                'message': 'You have successfully added {} coins!'.format(self.body.get('value')),
                                'email': self.user.email,
                                'body': {}
                            }
                            return make_response(jsonify(responseObject)), 200
                        raise ValueError({'status': 1, 'message': 'Wrong request!'}, 400)
                    elif (self.post_data.get('type') == 3): # Get directory tree
                        self.get_body()
                        if (len(self.body) == 0):
                            responseObject = {
                                'status': 0,
                                'type': 3,
                                'message': 'You have successfully got directory tree!',
                                'email': self.user.email,
                                'body': json.loads(self.user.data)
                            }
                            return make_response(jsonify(responseObject)), 200
                        raise ValueError({'status': 1, 'message': 'Wrong request!'}, 400)
                    elif (self.post_data.get('type') == 4): # Get coins
                        self.get_body()
                        if (len(self.body) == 0):
                            responseObject = {
                                'status': 0,
                                'type': 4,
                                'message': 'You have {} coins!'.format(self.user.coins),
                                'email': self.user.email,
                                'body': {
                                    'value': self.user.coins
                                }
                            }
                            return make_response(jsonify(responseObject)), 200
                        raise ValueError({'status': 1, 'message': 'Wrong request!'}, 400)
                    elif (self.post_data.get('type') == 5): # Add file or real directory
                        self.get_body()
                        if not isinstance(self.body.get('IsDir'), bool):
                            raise ValueError({'status': 1, 'message': 'You should specify IsDir flag!'}, 400)
                        if not isinstance(self.body.get("Size"), int):
                            raise ValueError({'status': 1, 'message': 'You should specify Size!'}, 400)
                        if (self.body.get("Size") > self.user.coins):
                            raise ValueError({'status': 1, 'message': 'You don\'t have enough coins!'}, 400)
                        item = {
                            "Size": self.body.get("Size"),
                            "Flag": True
                            }
                        if self.body.get('IsDir'): # Real directory
                            item["IsDir"] = True
                            item["Child"] = {}
                        else: # File
                            item["IsDir"] = False
                        return self.add_item(item, 5)
                    else:
                        raise ValueError({'status': 1, 'message': 'Wrong request!'}, 400)
                else:
                    return make_response(jsonify({'status': 1, 'message': resp})), 401
            else:
                return make_response(jsonify({'status': 1, 'message': 'Provide a valid auth token.'})), 403                

        except ValueError as responseObject:
            if (len(responseObject.args) == 2) and (isinstance(responseObject.args[1], int)):
                return make_response(jsonify(responseObject.args[0])), responseObject.args[1]
            else:
                return make_response(jsonify({'status': 1, 'message': 'Something went wrong'})), 400

        except KeyError:
            return make_response(jsonify({'status': 1, 'message': self.body.get('path') +  ': No such file or directory'})), 400
            
        except BadRequest:
            return make_response(jsonify({'status': 1, 'message': 'Request should be JSON type!'})), 400


class NodeUploadAPI(MethodView):
    """
    Node Upload Resource
    """

    def post(self):
        auth_header = request.headers.get('Authorization')
        if auth_header:
            try:
                auth_token = auth_header.split(" ")[1]
            except IndexError:
                responseObject = {
                    'status': 1,
                    'message': 'Bearer token malformed.'
                }
                return make_response(jsonify(responseObject)), 401
        else:
            auth_token = ''
        if auth_token:
            certificate = Сertificate.query.filter_by(token=auth_token).first()
            #if certificate:
            #    
            #else:
            return make_response(jsonify({'status': 1, 'message': resp})), 401
        else:
            return make_response(jsonify({'status': 1, 'message': 'Provide a valid auth token.'})), 401


class NodeDownloadAPI(MethodView):
    """
    Node Download Resource
    """
    pass


class NodeDeleteAPI(MethodView):
    """
    Node Delete Resource
    """
    pass


class AddNodeAPI(MethodView):
    """
    Add Node Resource
    """

    def post(self):
        try:
            # get the post data
            post_data = request.get_json()

            if not post_data: # Request isn't JSON type
                raise BadRequest
        
            if not post_data.get('ip_address'):
                raise ValueError({'status': 1, 'message': 'You must specify IP address!'}, 400)
            if not (len(post_data) == 1):
                raise ValueError({'status': 1, 'message': 'Too many arguments!'}, 400)
            
            # Get node
            ip_address_int = int(IPA(post_data.get('ip_address')))
            node = Node.query.filter_by(ip_address=ip_address_int).first()
            if not node: # if node doesn't exist
                new_node = Node(ip_address_int)
                
                if new_node:
                    db.session.add(new_node)
                    db.session.commit()

                    responseObject = {
                        'status': 0,
                        'message': 'Successfully added node!',
                    }
                    return make_response(jsonify(responseObject)), 200
                else:
                    raise ValueError({'status': 1, 'message': 'Can\'t load node\'s IP!'}, 400)
            else:
                raise ValueError({'status': 1, 'message': 'Node already exists!'}, 400)

        except BadRequest:
            return make_response(jsonify({'status': 1, 'message': 'Request should be JSON type!'})), 400

        except ValueError as responseObject:
            if (len(responseObject.args) == 2) and (isinstance(responseObject.args[1], int)):
                return make_response(jsonify(responseObject.args[0])), responseObject.args[1]
            else:
                return make_response(jsonify({'status': 1, 'message': 'Wrong IP format!'})), 400 

        except Exception as e:
            return make_response(jsonify({'status': 1, 'message': 'Some error occurred. Please try again.'})), 401


class DeleteNodeAPI(MethodView):
    """
    Delete Node Resource
    """

    def delete(self):
        try:
            # get the delete data
            delete_data = request.get_json()

            if not delete_data: # Request isn't JSON type
                raise BadRequest
        
            if not delete_data.get('ip_address'):
                raise ValueError({'status': 1, 'message': 'You must specify IP address!'}, 400)
            if not (len(delete_data) == 1):
                raise ValueError({'status': 1, 'message': 'Too many arguments!'}, 400)
            
            # Get node
            ip_address_int = int(IPA(delete_data.get('ip_address')))
            print(ip_address_int)
            node = Node.query.filter_by(ip_address=ip_address_int).first()
            if not node: # if node doesn't exist
                raise ValueError({'status': 1, 'message': 'Node doesn\'t exist!'}, 400)
            else:
                db.session.delete(node)
                db.session.commit()

                responseObject = {
                    'status': 0,
                    'message': 'Successfully deleted node!',
                }
                return make_response(jsonify(responseObject)), 200
        
        except ValueError as responseObject:
            if (len(responseObject.args) == 2) and (isinstance(responseObject.args[1], int)):
                return make_response(jsonify(responseObject.args[0])), responseObject.args[1]
            else:
                return make_response(jsonify({'status': 1, 'message': 'Wrong IP format!'})), 400 

        except Exception as e:
            return make_response(jsonify({'status': 1, 'message': 'Some error occurred. Please try again.'})), 401

        except BadRequest:
            return make_response(jsonify({'status': 1, 'message': 'Request should be JSON type!'})), 400

    
# define the API resources
registration_view = RegisterAPI.as_view('register_api')
login_view = LoginAPI.as_view('login_api')
user_view = UserAPI.as_view('user_api')
logout_view = LogoutAPI.as_view('logout_api')
request_view = RequestAPI.as_view('request_api')
upload_view = NodeUploadAPI.as_view('upload_api')
download_view = NodeDownloadAPI.as_view('download_api')
delete_view = NodeDeleteAPI.as_view('delete_api')
add_node_view = AddNodeAPI.as_view('add_node_api')
delete_node_view = DeleteNodeAPI.as_view('delete_node_api')

# add Rules for API Endpoints
auth_blueprint.add_url_rule(
    '/auth/register',
    view_func=registration_view,
    methods=['POST']
)
auth_blueprint.add_url_rule(
    '/auth/login',
    view_func=login_view,
    methods=['POST']
)
auth_blueprint.add_url_rule(
    '/auth/status',
    view_func=user_view,
    methods=['GET']
)
auth_blueprint.add_url_rule(
    '/auth/logout',
    view_func=logout_view,
    methods=['POST']
)
auth_blueprint.add_url_rule(
    '/auth/request',
    view_func=request_view,
    methods=['POST']
)
auth_blueprint.add_url_rule(
    '/auth/upload',
    view_func=upload_view,
    methods=['POST']
)
auth_blueprint.add_url_rule(
    '/auth/download',
    view_func=download_view,
    methods=['POST']
)
auth_blueprint.add_url_rule(
    '/auth/delete',
    view_func=delete_view,
    methods=['POST']
)
auth_blueprint.add_url_rule(
    '/auth/node/add',
    view_func=add_node_view,
    methods=['POST']
)
auth_blueprint.add_url_rule(
    '/auth/node/delete',
    view_func=delete_node_view,
    methods=['DELETE']
)