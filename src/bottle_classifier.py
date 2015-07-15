import dlib, sys, glob
from skimage import io
import datetime
import shutil
import time
import datetime
import os


bottle_training_data = 'helpers/bottles_dataset.xml'
img_list = 'fetch_images/image_urls.txt'
options = dlib.simple_object_detector_training_options()
options.add_left_right_image_flips = True

# options.detection_window_size = 8000
options.C = 4
options.epsilon = 0.01
options.num_threads = 8
options.be_verbose = True

dlib.train_simple_object_detector(bottle_training_data,"square_bottle_classifier.svm",options)

ts = time.time()
st = datetime.datetime.fromtimestamp(ts).strftime('%Y-%m-%d %H:%M:%S') + '/'

parentDir = 'models/'
if not os.path.exists(parentDir + st):
    os.makedirs(parentDir + st)
    shutil.copyfile("square_bottle_classifier.svm", parentDir + st + "square_bottle_classifier.svm")
    shutil.copyfile(bottle_training_data, parentDir + st + "bottles_dataset.xml")
    shutil.copyfile(img_list, parentDir + st + "image_urls.txt")

detector = dlib.simple_object_detector("square_bottle_classifier.svm")
win = dlib.image_window()
win.set_image(detector)

# print "\nTraining accuracy: ", dlib.test_simple_object_detector(bottle_training_data, "square_bottle_classifier.svm")

