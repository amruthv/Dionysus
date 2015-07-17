package antoine.photonoconfirm1;

import android.app.Activity;
import android.content.Context;
import android.content.pm.PackageManager;
import android.hardware.Camera;
import android.os.Bundle;
import android.widget.FrameLayout;

public class LastDitchActivity extends Activity {

    private Camera mCamera;
    private LastDitchSurfaceView  mPreview;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.last_ditch_layout);

        mCamera = getCameraInstance();
        mPreview = new LastDitchSurfaceView(this, mCamera);

        FrameLayout preview = (FrameLayout)findViewById(R.id.last_ditch_frame);
        preview.addView(mPreview);
    }

    /** Check if this device has a camera only if not specified in the manifest */
    public boolean checkCameraHardware(Context context) {
        if (context.getPackageManager().hasSystemFeature(PackageManager.FEATURE_CAMERA)){
            // this device has a camera
            return true;
        } else {
            // no camera on this device
            return false;
        }
    }

    /** A safe way to get an instance of the Camera object. */
    public static Camera getCameraInstance(){
        Camera c = null;
        try {
            c = Camera.open(); // attempt to get a Camera instance
        } catch (Exception e) {
            // Camera is not available (in use or does not exist)
        }
        return c; // returns null if camera is unavailable
    }

    /**Check if the device has flash*/
    public boolean checkFlash(Context context){
        if(context.getPackageManager().hasSystemFeature(PackageManager.FEATURE_CAMERA_FLASH)){
            //the device has flash
            return true;
        }else{
            //no flash
            return false;
        }

    }

    @Override
    protected void onDestroy() {
        // TODO Auto-generated method stub
        super.onDestroy();
        releaseCamera();
    }

    @Override
    protected void onPause() {
        // TODO Auto-generated method stub
        super.onPause();
        releaseCamera();
    }

    @Override
    protected void onResume() {
        // TODO Auto-generated method stub
        super.onResume();

        //Test if i have to put all this code like in onCreate
        if(mCamera!=null){
            return;
        }
        mCamera=getCameraInstance();

        if(mPreview!=null){
            return;
        }
        mPreview = new LastDitchSurfaceView(this, mCamera);
        FrameLayout preview = (FrameLayout)findViewById(R.id.last_ditch_frame);
        preview.addView(mPreview);
    }

    private void releaseCamera(){
        if (mCamera != null){
            mCamera.release();        // release the camera for other applications
            mCamera = null;
        }
    }
}