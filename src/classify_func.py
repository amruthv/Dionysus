import os
import dlib
import sys
from skimage import io
import numpy as np

detector = dlib.simple_object_detector("square_bottle_classifier.svm")

assorted_dir = '../test_images/'

items =os.listdir(assorted_dir)

def classify(dir):
  for dirr, _, files in os.walk(assorted_dir):
    sorted_files = sorted(files)
    for f in files:
      if f.endswith('.jpeg'):
        im = io.imread(dirr + f)
        dets = detector(im)
        print f + ' had: ', len(dets)
if __name__ == '__main__':
  classify(assorted_dir)