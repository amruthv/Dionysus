#include <dlib/svm_threaded.h>
#include <dlib/gui_widgets.h>
#include <dlib/image_processing.h>
#include <dlib/data_io.h>

#include <iostream>
#include <fstream>
#include <string>

#include <time.h>
#include <stdlib.h>
#include <stdio.h>
#include <sys/stat.h>


using namespace std;
using namespace dlib;



int main(int argc, char** argv)
{  

    try
    {
        if (argc != 2)
        {
            cout << "Give the path to the training images directory as the argument to this program." << endl;
            return 0;
        }
        const std::string images_directory = argv[1];
        
        dlib::array<array2d<unsigned char> > images_test;
        std::vector<std::vector<rectangle> > face_boxes_test;

        load_image_dataset(images_test, face_boxes_test, images_directory+"/test_images.xml");
        cout << "num testing images:  " << images_test.size() << endl;

        typedef scan_fhog_pyramid<pyramid_down<6> > image_scanner_type; 
        object_detector<image_scanner_type> detector;
        deserialize("bottle_classifier.svm") >> detector;

        cout << "testing results:  " << test_object_detection_function(detector, images_test, face_boxes_test) << endl;        
        image_window hogwin(draw_fhog(detector), "Learned fHOG detector");
        image_window win; 
        for (unsigned long i = 0; i < images_test.size(); ++i) {
            std::vector<rectangle> dets = detector(images_test[i]);
            
            cout << "number of detections" << dets.size() << endl;
            if (dets.size() != 0) {
                win.clear_overlay();
                win.set_image(images_test[i]);
                win.add_overlay(dets, rgb_pixel(255,0,0));
                cout << "Hit enter to process the next image..." << endl;
                cin.get();
            }
        }
    } catch (exception& e) {
        cout << "\nexception thrown!" << endl;
        cout << e.what() << endl;
    }
}