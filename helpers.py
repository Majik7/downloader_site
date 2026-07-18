import uuid, os, subprocess

def download_video(url: str):
    command = 'yt-dlp --ignore-errors --no-warnings --force-overwrites --merge-output-format mp4 -o "video.mp4" '

    res = subprocess.run(command.split()+[url], capture_output=True, text=True)
    out = res.stdout