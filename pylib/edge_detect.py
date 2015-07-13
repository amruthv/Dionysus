#! /usr/local/bin/python
import sys
import cv2
import numpy as np
from matplotlib import pyplot as plt

DEFAULT_IMAGE = "../res/image1.jpg"

if __name__ == "__main__":
  imagename = sys.argv[1] if len(sys.argv) > 1 else DEFAULT_IMAGE

  img = cv2.imread(imagename)
  img = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
  img = cv2.GaussianBlur(img, (5, 5), 0)
  edges = cv2.Canny(img, 100, 200)

  plt.subplot(121),plt.imshow(img,cmap = 'gray')
  plt.title('Original Image'), plt.xticks([]), plt.yticks([])
  plt.subplot(122),plt.imshow(edges,cmap = 'gray')
  plt.title('Edge Image'), plt.xticks([]), plt.yticks([])

  plt.show()
