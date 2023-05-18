import random
import string
from flask import Flask

app = Flask(__name__)

@app.route('/password')
def generate_password():
    password_list = []
    for i in range(100):
        password = ''.join(random.choices(string.ascii_letters + string.digits + '#@$', k=12))
        password_list.append(password)
    return str(password_list)

if __name__ == '__main__':
    app.run(port=8080)
