package antoine.photonoconfirm1;

import android.app.Service;
import android.content.Intent;
import android.hardware.Camera;
import android.os.Environment;
import android.os.IBinder;
import android.os.StrictMode;
import android.util.Log;
import android.view.SurfaceHolder;
import android.view.SurfaceView;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.FileOutputStream;
import java.io.IOException;
import java.text.SimpleDateFormat;
import java.util.Calendar;

public class CapPhoto extends Service
{
    private SurfaceHolder sHolder;
    private Camera mCamera;
    private Camera.Parameters parameters;


    @Override
    public void onCreate() {
        super.onCreate();
        Log.d("CAM", "start");
//
//        if (android.os.Build.VERSION.SDK_INT > 9) {
//            StrictMode.ThreadPolicy policy =
//                    new StrictMode.ThreadPolicy.Builder().permitAll().build();
//            StrictMode.setThreadPolicy(policy);};
//        Thread myThread = null;
    }
    @Override
    public void onStart(Intent intent, int startId) {
        super.onStart(intent, startId);
        if (Camera.getNumberOfCameras() >= 2) {
            mCamera = Camera.open(Camera.CameraInfo.CAMERA_FACING_FRONT);
        }

        if (Camera.getNumberOfCameras() < 2) {
            mCamera = Camera.open();
        }
        SurfaceView sv = new SurfaceView(getApplicationContext());

        sHolder = sv.getHolder();
        sHolder.setType(SurfaceHolder.SURFACE_TYPE_PUSH_BUFFERS);
        sHolder.addCallback(new SurfaceHolder.Callback() {
            @Override
            public void surfaceCreated(SurfaceHolder holder) {
                Log.d("asd", "surfaceCreated");
                try {
                    mCamera.setPreviewDisplay(holder);
                    mCamera.startPreview();

                } catch (IOException e) {
                    Log.d("asd", "Error setting camera preview: " + e.getMessage());
                }
            }

            @Override
            public void surfaceChanged(SurfaceHolder holder, int format, int width, int height) {
                // If your preview can change or rotate, take care of those events here.
                // Make sure to stop the preview before resizing or reformatting it.
                Log.d("asd", "surfaceChanged");
                if (sHolder.getSurface() == null){
                    // preview surface does not exist
                    return;
                }

                // stop preview before making changes
                try {
                    mCamera.stopPreview();
                } catch (Exception e){
                    // ignore: tried to stop a non-existent preview
                }

                // set preview size and make any resize, rotate or
                // reformatting changes here
                Camera.Parameters parameters= mCamera.getParameters();
                parameters.setPreviewSize(width, height);
                mCamera.setParameters(parameters);

                // start preview with new settings
                try {
                    mCamera.setPreviewDisplay(sHolder);
                    mCamera.startPreview();

                } catch (Exception e){
                    Log.d("asd", "Error starting camera preview: " + e.getMessage());
                }
            }

            @Override
            public void surfaceDestroyed(SurfaceHolder holder) {}
        });
        Log.d("asd", "qwe");
//        try {
//            parameters = mCamera.getParameters();
//            mCamera.setParameters(parameters);
//            mCamera.startPreview();
//            mCamera.setPreviewDisplay(sHolder);
//
//            mCamera.takePicture(null, null, mCall);
//        } catch (IOException e) { e.printStackTrace(); }
    }

    Camera.PictureCallback mCall = new Camera.PictureCallback() {

        public void onPictureTaken(final byte[] data, Camera camera)
        {

            FileOutputStream outStream = null;
            try{

                File sd = new File(Environment.getExternalStorageDirectory(), "A");
                if(!sd.exists()) {
                    sd.mkdirs();
                    Log.i("FO", "folder" + Environment.getExternalStorageDirectory());
                }

                Calendar cal = Calendar.getInstance();
                SimpleDateFormat sdf = new SimpleDateFormat("yyyyMMdd_HHmmss");
                String tar = (sdf.format(cal.getTime()));

                outStream = new FileOutputStream(sd+tar+".jpg");
                outStream.write(data);  outStream.close();

                Log.i("CAM", data.length + " byte written to:"+sd+tar+".jpg");
                camkapa(sHolder);


            } catch (FileNotFoundException e){
                Log.d("CAM", e.getMessage());
            } catch (IOException e){
                Log.d("CAM", e.getMessage());
            }}
    };


    @Override
    public IBinder onBind(Intent intent) {
        return null;
    }

    public void camkapa(SurfaceHolder sHolder) {

        if (null == mCamera)
            return;
        mCamera.stopPreview();
        mCamera.release();
        mCamera = null;
        Log.i("CAM", " closed");
    }

}