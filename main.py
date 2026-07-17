from flask import Flask, request, render_template, redirect, url_for
from markupsafe import escape
import os

UPLOAD_DIR = "uploads"

app = Flask(__name__)

@app.get("/")
def welcome():
    return "<p>Welcome</p>"

@app.get('/hello/')
@app.get('/hello/<name>')
def hello(name=None):
    return render_template('hello.html', person=name)

@app.post("/upload")
def upload_txt():
    f = request.files["file"]
    os.makedirs(UPLOAD_DIR, exist_ok=True)
    f.save(os.path.join(UPLOAD_DIR, "uploaded.txt"))

    return "Uploaded"