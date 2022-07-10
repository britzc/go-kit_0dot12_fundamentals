import os
import shutil 
import subprocess

src = '.course/start/module05' 
dest = 'current'

if os.path.exists("current") and os.path.isdir("current"):
    shutil.rmtree("current")   
    
shutil.copytree(src, dest) 

p = subprocess.Popen(["go", "mod", "init"], cwd="current")
p.wait()

p = subprocess.Popen(["go", "mod", "tidy"], cwd="current")
p.wait()

