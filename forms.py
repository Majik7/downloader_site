from flask_wtf import FlaskForm
from wtforms import StringField, SubmitField
from wtforms.validators import DataRequired

class URLForm(FlaskForm):
    url = StringField("Enter video URL", validators=[DataRequired()])
    submit = SubmitField("Submit")