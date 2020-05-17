# project/server/auth/views.py

from flask import Blueprint, request, make_response, jsonify
from flask.views import MethodView

from project.server import bcrypt, db
from project.server.models import User, BlacklistToken, Сertificate

import json
from functools import reduce  # forward compatibility for Python 3
import operator
from werkzeug.exceptions import BadRequest

auth_blueprint = Blueprint('auth', __name__)

class RegisterAPI(MethodView):
    """
    User Registration Resource
    """

    def post(self):
        # get the post data
        post_data = request.get_json()
        # check if user already exists
        user = User.query.filter_by(email=post_data.get('email')).first()
        if not user:
            try:
                user = User(
                    email=post_data.get('email'),
                    password=post_data.get('password')
                )

                # insert the user
                db.session.add(user)
                db.session.commit()
                # generate the auth token
                auth_token = user.encode_auth_token(user.id)
                responseObject = {
                    'status': 'success',
                    'message': 'Successfully registered.',
                    'auth_token': auth_token.decode()
                }
                return make_response(jsonify(responseObject)), 201
            except Exception as e:
                responseObject = {
                    'status': 'fail',
                    'message': 'Some error occurred. Please try again.'
                }
                return make_response(jsonify(responseObject)), 401
        else:
            responseObject = {
                'status': 'fail',
                'message': 'User already exists. Please Log in.',
            }
            return make_response(jsonify(responseObject)), 202

class LoginAPI(MethodView):
    """
    User Login Resource
    """
    def post(self):
        # get the post data
        post_data = request.get_json()
        try:
            # fetch the user data
            user = User.query.filter_by(
                email=post_data.get('email')
            ).first()
            if user and bcrypt.check_password_hash(
                user.password, post_data.get('password')
            ):
                auth_token = user.encode_auth_token(user.id)
                if auth_token:
                    responseObject = {
                        'status': 'success',
                        'message': 'Successfully logged in.',
                        'auth_token': auth_token.decode()
                    }
                    return make_response(jsonify(responseObject)), 200
            elif user and not bcrypt.check_password_hash( # If incorrect password
                user.password, post_data.get('password')
            ):
                responseObject = {
                    'status': 'fail',
                    'message': 'Incorrect password.'
                }
                return make_response(jsonify(responseObject)), 404
            else:
                responseObject = {
                    'status': 'fail',
                    'message': 'User does not exist.'
                }
                return make_response(jsonify(responseObject)), 404
        except Exception as e:
            print(e)
            responseObject = {
                'status': 'fail',
                'message': 'Try again'
            }
            return make_response(jsonify(responseObject)), 500

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
                    'status': 'fail',
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
                    'status': 'success',
                    'data': {
                        'user_id': user.id,
                        'email': user.email,
                        'admin': user.admin,
                        'registered_on': user.registered_on
                    }
                }
                return make_response(jsonify(responseObject)), 200
            responseObject = {
                'status': 'fail',
                'message': resp
            }
            return make_response(jsonify(responseObject)), 401
        else:
            responseObject = {
                'status': 'fail',
                'message': 'Provide a valid auth token.'
            }
            return make_response(jsonify(responseObject)), 401


class LogoutAPI(MethodView):
    """
    Logout Resource
    """
    def post(self):
        # get auth token
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
                    responseObject = {
                        'status': 'success',
                        'message': 'Successfully logged out.'
                    }
                    return make_response(jsonify(responseObject)), 200
                except Exception as e:
                    responseObject = {
                        'status': 'fail',
                        'message': e
                    }
                    return make_response(jsonify(responseObject)), 200
            else:
                responseObject = {
                    'status': 'fail',
                    'message': resp
                }
                return make_response(jsonify(responseObject)), 401
        else:
            responseObject = {
                'status': 'fail',
                'message': 'Provide a valid auth token.'
            }
            return make_response(jsonify(responseObject)), 403


class RequestAPI(MethodView):
    """
    User Update Resource
    """

    def get_body(self):
        if (len(self.post_data) > 5):
            raise ValueError({'status': 1, 'message': 'Too many fields!'}, 400)
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


    def set_file_by_path(self, root, items, value):
        """Set a file in a nested object in root by item sequence."""
        if value in self.get_by_path(root, items):
            raise ValueError({'status': 1, 'message': 'File already exists!'}, 400)
        self.get_by_path(root, items)[value] = {
                                        "Name": value,
                                        "Size": 0,
                                        "IsDir": False,
                                        "Flag": False
                                        }
    

    def set_directory_by_path(self, root, items, value):
        """Set a directory in a nested object in root by item sequence."""
        if value in self.get_by_path(root, items):
            raise ValueError({'status': 1, 'message': 'Directory already exists!'}, 400)
        self.get_by_path(root, items)[value] = {
                                        "Name": value,
                                        "Size": 0,
                                        "IsDir": True,
                                        "Flag": False,
                                        "Child": {}
                                        }

    def delete_by_path(self, root, items, value):
        """Delete a directory or file in a nested object in root by item sequence."""
        del self.get_by_path(root, items)[value]

    def post(self):
        try:
            # get the post data
            self.post_data = request.get_json()

            if not self.post_data: # Request isn't JSON type
                raise BadRequest
            
            if (self.post_data.get('type') == 0): # Add directory
                self.get_body()
                if not self.body.get('name'):
                    raise ValueError({'status': 1, 'message': 'You should specify name in "Add directory" method!'}, 400)
                if not self.body.get('path'):
                    raise ValueError({'status': 1, 'message': 'You should specify path in "Add directory" method!'}, 400)
                if not (len(self.body) == 2):
                    raise ValueError({'status': 1, 'message': 'Too many arguments in "Add directory" method!'}, 400)
                
                path = self.body.get('path') + self.body.get('name')
                abs_path = ["Child"]
                if (path[0] != '/'):
                    raise ValueError({'status': 1, 'message': 'You must always start your path from "/" symbol!'}, 400)
                if (len(path) == 1):
                    raise ValueError({'status': 1, 'message': 'You must specify directory!'}, 400)
                data = json.loads(self.user.data)
                if (path[-1] == '/'):
                    initial_path = path.split('/')[1:-1]
                else:
                    initial_path = path.split('/')[1:]
                
                abs_path = ["Child"]
                for i in initial_path[:-1]:
                    abs_path.append(i)
                    abs_path.append("Child")
                
                self.set_directory_by_path(data, abs_path, initial_path[-1])
                self.user.data = json.dumps(data)
                db.session.commit()
                responseObject = {
                    'status': 0,
                    'type': 0,
                    'message': 'You have successfully added new directory!',
                    'email': self.user.email,
                    'body': {}
                }
                return make_response(jsonify(responseObject)), 200
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
            elif (self.post_data.get('type') == 5): # Add file
                self.get_body()
                if not self.body.get('name'):
                    raise ValueError({'status': 1, 'message': 'You should specify name in "Add file" method!'}, 400)
                if not self.body.get('path'):
                    raise ValueError({'status': 1, 'message': 'You should specify path in "Add file" method!'}, 400)
                if not (len(self.body) == 2):
                    raise ValueError({'status': 1, 'message': 'Too many arguments in "Add file" method!'}, 400)
                
                path = self.body.get('path') + self.body.get('name')
                abs_path = ["Child"]
                if (path[0] != '/'):
                    raise ValueError({'status': 1, 'message': 'You must always start your path from "/" symbol!'}, 400)
                if (len(path) == 1):
                    raise ValueError({'status': 1, 'message': 'You must specify file!'}, 400)
                data = json.loads(self.user.data)
                if (path[-1] == '/'):
                    initial_path = path.split('/')[1:-1]
                else:
                    initial_path = path.split('/')[1:]
                
                abs_path = ["Child"]
                for i in initial_path[:-1]:
                    abs_path.append(i)
                    abs_path.append("Child")
                
                self.set_file_by_path(data, abs_path, initial_path[-1])
                self.user.data = json.dumps(data)
                db.session.commit()
                responseObject = {
                    'status': 0,
                    'type': 0,
                    'message': 'You have successfully added new file!',
                    'email': self.user.email,
                    'body': {}
                }
                return make_response(jsonify(responseObject)), 200
            else:
                raise ValueError({'status': 1, 'message': 'Wrong request!'}, 400)

        except ValueError as responseObject:
            return make_response(jsonify(responseObject.args[0])), responseObject.args[1]

        except KeyError:
            return make_response(jsonify({'status': 1, 'message': self.body.get('path') +  ': No such file or directory'})), 400

        except BadRequest:
            return make_response(jsonify({'status': 1, 'message': 'Request should be JSON type!'})), 400


# define the API resources
registration_view = RegisterAPI.as_view('register_api')
login_view = LoginAPI.as_view('login_api')
user_view = UserAPI.as_view('user_api')
logout_view = LogoutAPI.as_view('logout_api')
request_view = RequestAPI.as_view('request_api')

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
    methods=['POST', 'GET']
)
