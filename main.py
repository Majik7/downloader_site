from flask import Flask, request, render_template, redirect, url_for, send_file, after_this_request
from markupsafe import escape
import os, uuid
import subprocess
from flask_bootstrap import Bootstrap5
from flask_wtf import FlaskForm, CSRFProtect
from wtforms import StringField, SubmitField
from wtforms.validators import DataRequired
from forms import URLForm
from dotenv import load_dotenv

UPLOAD_DIR = "uploads"
DOWNLOAD_DIR = "downloads"

app = Flask(__name__)
app.secret_key = os.environ.get("SECRET_KEY")

bootstrap = Bootstrap5(app)
csrf = CSRFProtect(app)

@app.get("/")
def welcome():
    return "<p>Welcome</p>"

@app.get('/hello/')
@app.get('/hello/<name>')
def hello(name=None):
    return render_template('hello.html', person=name)

@app.route("/app", methods=["GET", "POST"])
def main_form():
    form = URLForm()

    if form.validate_on_submit():
        url = form.url.data
        return redirect(url_for("download_video", url=url))

    return render_template("main_form.html", form=form)

@app.get("/shelltest")
def shelltest():
    res = subprocess.run(["python", "--version"], capture_output=True, text=True)
    out = res.stdout

    app.logger.debug("SHELL COMMAND CHAL GAYA OYEEE")

    return render_template("shelltest.html", out=out)

@app.route('/download-video', methods=['GET', 'POST'])
def download_video():
    for f in os.listdir(DOWNLOAD_DIR):
        try:
            os.remove(os.path.join(DOWNLOAD_DIR, f))
            app.logger.debug(f"Deleting file {f}")
        except OSError:
            pass
    
    url = request.args.get('url') or request.form.get('url')
    if not url:
        return "No URL provided", 400

    os.makedirs(DOWNLOAD_DIR, exist_ok=True)
    filepath = os.path.join(DOWNLOAD_DIR, f"{uuid.uuid4()}.mp4")

    result = subprocess.run(
        ["yt-dlp", "-o", filepath, "--merge-output-format", "mp4", "-f",
         "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best", url],
        capture_output=True, text=True
    )

    if result.returncode != 0:
        app.logger.error(f"yt-dlp failed: {result.stderr}")
        return "Download failed", 500

    response = send_file(filepath, as_attachment=True, download_name="video.mp4")

    return response