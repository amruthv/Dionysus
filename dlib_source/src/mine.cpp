// The contents of this file are in the public domain. See LICENSE_FOR_EXAMPLE_PROGRAMS.txt
/*

    This example program shows how you can use dlib to make an object detector
    for things like faces, pedestrians, and any other semi-rigid object.  In
    particular, we go though the steps to train the kind of sliding window
    object detector first published by Dalal and Triggs in 2005 in the paper
    Histograms of Oriented Gradients for Human Detection.  

    Note that this program executes fastest when compiled with at least SSE2
    instructions enabled.  So if you are using a PC with an Intel or AMD chip
    then you should enable at least SSE2 instructions.  If you are using cmake
    to compile this program you can enable them by using one of the following
    commands when you create the build project:
        cmake path_to_dlib_root/examples -DUSE_SSE2_INSTRUCTIONS=ON
        cmake path_to_dlib_root/examples -DUSE_SSE4_INSTRUCTIONS=ON
        cmake path_to_dlib_root/examples -DUSE_AVX_INSTRUCTIONS=ON
    This will set the appropriate compiler options for GCC, clang, Visual
    Studio, or the Intel compiler.  If you are using another compiler then you
    need to consult your compiler's manual to determine how to enable these
    instructions.  Note that AVX is the fastest but requires a CPU from at least
    2011.  SSE4 is the next fastest and is supported by most current machines.  

*/


#include <dlib/svm_threaded.h>
#include <dlib/image_processing.h>
#include <dlib/data_io.h>

#include <iostream>
#include <fstream>


using namespace std;
using namespace dlib;

// ----------------------------------------------------------------------------------------

int main(int argc, char** argv)
{  

    try
    {
        // In this example we are going to train a face detector based on the
        // small faces dataset in the examples/faces directory.  So the first
        // thing we do is load that dataset.  This means you need to supply the
        // path to this faces folder as a command line argument so we will know
        // where it is.
        if (argc != 2)
        {
            cout << "Give the path to the examples/faces directory as the argument to this" << endl;
            cout << "program.  For example, if you are in the examples folder then execute " << endl;
            cout << "this program by running: " << endl;
            cout << "   ./fhog_object_detector_ex faces" << endl;
            cout << endl;
            return 0;
        }
        const std::string faces_directory = argv[1];
        dlib::array<array2d<unsigned char> > images_train, images_test;
        std::vector<std::vector<rectangle> > face_boxes_train, face_boxes_test;

        load_image_dataset(images_train, face_boxes_train, faces_directory+"/training.xml");
        load_image_dataset(images_test, face_boxes_test, faces_directory+"/testing.xml");

        upsample_image_dataset<pyramid_down<2> >(images_train, face_boxes_train);
        upsample_image_dataset<pyramid_down<2> >(images_test,  face_boxes_test);
        add_image_left_right_flips(images_train, face_boxes_train);
        cout << "num training images: " << images_train.size() << endl;
        cout << "num testing images:  " << images_test.size() << endl;


        typedef scan_fhog_pyramid<pyramid_down<6> > image_scanner_type; 
        image_scanner_type scanner;
        // The sliding window detector will be 80 pixels wide and 80 pixels tall.
        scanner.set_detection_window_size(80, 80); 
        structural_object_detection_trainer<image_scanner_type> trainer(scanner);
        // Set this to the number of processing cores on your machine.
        trainer.set_num_threads(4);  
        trainer.set_c(1);
        trainer.be_verbose();
        trainer.set_epsilon(0.01);


        object_detector<image_scanner_type> detector = trainer.train(images_train, face_boxes_train);

        cout << "training results: " << test_object_detection_function(detector, images_train, face_boxes_train) << endl;
        // However, to get an idea if it really worked without overfitting we need to run
        // it on images it wasn't trained on.  The next line does this.  Happily, we see
        // that the object detector works perfectly on the testing images.
        cout << "testing results:  " << test_object_detection_function(detector, images_test, face_boxes_test) << endl;


        // If you have read any papers that use HOG you have probably seen the nice looking
        // "sticks" visualization of a learned HOG detector.  This next line creates a
        // window with such a visualization of our detector.  It should look somewhat like
        // a face.
        serialize("face_detector.svm") << detector;
    } catch (exception& e) {
        cout << "\nexception thrown!" << endl;
        cout << e.what() << endl;
    }
}

// ----------------------------------------------------------------------------------------

