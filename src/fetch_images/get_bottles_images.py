import urllib
import os

this_dir = os.path.dirname(os.path.realpath(__file__))
numImages = 55
prefix = this_dir + '/../../dataset/'
f = open(this_dir + '/image_urls.txt')
print prefix


blacklist = ['beer-bottles.jpg']
count = 0
for line in f:
    print line
    if count > numImages:
        break
    if line.strip() not in blacklist:
        try:
            urllib.urlretrieve(line, prefix + line.strip().split('/')[-1])
        except:
            continue
        count += 1
f.close()

# clean up empty images
to_delete = []
for i in os.listdir(prefix):
    if os.path.getsize(prefix + i) == 2051:
        to_delete.append(prefix+ i)
for toDelete in to_delete:
  os.remove(toDelete)

