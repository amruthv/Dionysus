import os
import dlib
import sys
from skimage import io
import numpy as np

test_images_dir = '../test_images/'
default_svm_param_file = "square_bottle_classifier.svm"

items =os.listdir(test_images_dir)

def classify(dir, svm_param_file):
    detector = dlib.simple_object_detector(svm_param_file)
    for dirr, _, files in os.walk(test_images_dir):
        sorted_files = sorted(files)
        for f in files:
            if f.endswith('.jpeg'):
                im = io.imread(dirr + f)
                dets = detector(im)
                print f + ' had: ', len(dets)

if __name__ == '__main__':
    print sys.argv
    if len(sys.argv) == 2:
        classify(test_images_dir, sys.argv[1])
    else:
        #use default svm file
        classify(test_images_dir, default_svm_param_file)