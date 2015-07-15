#! /bin/bash

cd src/build
cmake .. && cmake --build . --config Release
cp train_bottle_detector ../../../src/ 
cp test_bottle_detector ../../../src/
