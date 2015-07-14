import os
import dlib
import sys
from skimage import io
import numpy as np

test_images_dir = '../test_images/'
default_svm_param_file = "square_bottle_classifier.svm"
rubric = '../test_images/rubric'

def classify(dir, svm_param_file):
  detMap = {}
  detector = dlib.simple_object_detector(svm_param_file)
  for dirr, _, files in os.walk(assorted_dir):
    for f in files:
      if f.endswith('.jpg'):
        im = io.imread(dirr + f)
        dets = detector(im)

        # split filename take first piece
        # match 'filename' with detected count
        detMap[f.split('.')[0]] = len(dets)
        print f + ' had: ', len(dets)

  print "score = " + score(detMap)

def score(detMap):
    sc = 0
    with open(rubric) as f:
        content = f.readlines()

    for line in content:
        words = line.split(' ')
        if int(words[0]) in detMap:
            sc += getScore(words[1:], detMap[int(words[0])])

    return float(sc) / float(len(detMap))

def getScore(words, count):
    rubric = {}
    for word in words:
        scorepair = word.split(":")
        rubric[int(scorepair[0])] = int(scorepair[1])

    if count in rubric:
        return rubric[count]

    return rubric[-1]

if __name__ == '__main__':
    print sys.argv
    if len(sys.argv) == 2:
        classify(test_images_dir, sys.argv[1])
    else:
        #use default svm file
        classify(test_images_dir, default_svm_param_file)

# detector('convert_dir/image3.jpg')
