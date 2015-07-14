import dlib, sys, glob
from skimage import io

bottle_training = 'helpers/bottles_dataset.xml'
options = dlib.simple_object_detector_training_options()
options.add_left_right_image_flips = True

options.C = 4
options.epsilon = 0.01
options.num_threads = 8
options.be_verbose = True

dlib.train_simple_object_detector(bottle_training,"square_bottle_classifier.svm",options)

# print "\nTraining accuracy: ", dlib.test_simple_object_detector(cups_training,"cupdetector.svm")

