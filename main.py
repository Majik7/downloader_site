from flask import Flask, request, render_template, redirect, url_for, send_file
from markupsafe import escape
import os
import subprocess

UPLOAD_DIR = "uploads"
DOWNLAOD_DIR = "downloads"

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
    url = request.args.get("url", None)
    f = request.files["file"]
    os.makedirs(UPLOAD_DIR, exist_ok=True)
    f.save(os.path.join(UPLOAD_DIR, "uploaded.txt"))

    return "Uploaded"

@app.get("/shelltest")
def shelltest():
    res = subprocess.run(["python", "--version"], capture_output=True, text=True)
    out = res.stdout

    app.logger.debug("SHELL COMMAND CHAL GAYA OYEEE")

    return render_template("shelltest.html", out=out)