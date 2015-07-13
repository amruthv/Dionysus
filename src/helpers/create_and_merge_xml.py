import subprocess

#create the new xml
bashCommand = "./imglab -c temp.xml ../../dataset/"
process = subprocess.Popen(bashCommand.split(), stdout=subprocess.PIPE)
output = process.communicate()[0]

# create an xml that merges the new and the old
bashCommand = "./imglab --add bottles_dataset.xml temp.xml"
process = subprocess.Popen(bashCommand.split(), stdout=subprocess.PIPE)
output = process.communicate()[0]

# move the merge.xml to override the old bottles_dataset.xml
bashCommand = "mv merged.xml bottles_dataset.xml"
process = subprocess.Popen(bashCommand.split(), stdout=subprocess.PIPE)
output = process.communicate()[0]

# delete the temp.xml
bashCommand = "rm temp.xml" 
process = subprocess.Popen(bashCommand.split(), stdout=subprocess.PIPE)
output = process.communicate()[0]
