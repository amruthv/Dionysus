import subprocess
bashCommand = "./imglab -c bottles_dataset.xml ../../dataset/"
process = subprocess.Popen(bashCommand.split(), stdout=subprocess.PIPE)
output = process.communicate()[0]
