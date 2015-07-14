import os
import dlib
import sys
from skimage import io
import numpy as np

# print "\nTest1 accuracy: ", dlib.test_simple_object_detector('/home/jyotiska/Dropbox/Computer Vision/cupdataset_2_test.xml',"cupdetector_2.svm")
# print "\nTraining accuracy: ", dlib.test_simple_object_detector('/home/jyotiska/Dropbox/Computer Vision/cupdataset_3.xml',"cupdetector_3.svm")

detector = dlib.simple_object_detector("square_bottle_classifier.svm")

# win_det = dlib.image_window()
# win_det.set_image(detector)

# win = dlib.image_window()
# test_dir = '/home/jyotiska/Dropbox/Computer Vision/Cups_test'
# convert_dir = '/home/jyotiska/Dropbox/Computer Vision/Cups_test_convert'
assorted_dir = '../test_images/'

items =os.listdir(assorted_dir)

def classify(dir):
  for dirr, _, files in os.walk(assorted_dir):
    for f in files:
      if f.endswith('.jpg'):
        im = io.imread(dirr + f)
        dets = detector(im)
        print f + ' had: ', len(dets)
if __name__ == '__main__':
  classify(assorted_dir)


# detector('convert_dir/image3.jpg')