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
        if (argc != 3)
        {
            cout << "Give the path to the examples/faces directory as the argument to this" << endl;
            cout << "program.  For example, if you are in the examples folder then execute " << endl;
            cout << "this program by running: " << endl;
            cout << "   ./fhog_object_detector_ex faces" << endl;
            cout << endl;
            return 0;
        }
        const std::string faces_directory = argv[1];
        const std::string item_to_classify = argv[2];
        
        dlib::array<array2d<unsigned char> > images_train, images_test;
        std::vector<std::vector<rectangle> > face_boxes_train, face_boxes_test;

        if (item_to_classify == "bottle") {
            load_image_dataset(images_train, face_boxes_train, faces_directory+"/bottles_dataset.xml");
            load_image_dataset(images_test, face_boxes_test, faces_directory+"/test_images.xml");
        } else {
            load_image_dataset(images_train, face_boxes_train, faces_directory+"/cans_dataset.xml");
            load_image_dataset(images_test, face_boxes_test, faces_directory+"/test_can_images.xml");
        }

        // upsample_image_dataset<pyramid_down<2> >(images_train, face_boxes_train);
        // upsample_image_dataset<pyramid_down<2> >(images_test,  face_boxes_test);

        add_image_left_right_flips(images_train, face_boxes_train);
        cout << "num training images: " << images_train.size() << endl;
        cout << "num testing images:  " << images_test.size() << endl;

        typedef scan_fhog_pyramid<pyramid_down<6> > image_scanner_type; 
        image_scanner_type scanner;
        
        scanner.set_detection_window_size(46, 138); 
        structural_object_detection_trainer<image_scanner_type> trainer(scanner);
        
        trainer.set_num_threads(8);
        
        trainer.set_c(10);

        trainer.set_epsilon(0.01);
        
        trainer.be_verbose();
                
        
        remove_unobtainable_rectangles(trainer, images_train, face_boxes_train);

        object_detector<image_scanner_type> detector = trainer.train(images_train, face_boxes_train);

        cout << "training results: " << test_object_detection_function(detector, images_train, face_boxes_train) << endl;

        time_t rawtime;
        struct tm * timeinfo;
        char buffer [80];
        time (&rawtime);
        timeinfo = localtime(&rawtime);
        strftime(buffer,80,"%F %R:%S",timeinfo);
        
        std::string timestamp_dir = "models/" + string(buffer);
        
        mkdir(timestamp_dir.c_str(), 0777);
        serialize(timestamp_dir + "/" + item_to_classify + "_classifier.svm") << detector;
        serialize(item_to_classify + "_classifier.svm") << detector;

        std::ifstream src1("helpers/" + item_to_classify + "_dataset.xml", std::ios::binary);
        std::ofstream dst1(timestamp_dir + item_to_classify +  "_dataset.xml", std::ios::binary);
        dst1 << src1.rdbuf();

        std::ifstream src2("fetch_images/image_urls.txt", std::ios::binary);
        std::ofstream dst2(timestamp_dir + "/image_urls.txt", std::ios::binary);
        dst2 << src2.rdbuf();
    } catch (exception& e) {
        cout << "\nexception thrown!" << endl;
        cout << e.what() << endl;
    }
}
