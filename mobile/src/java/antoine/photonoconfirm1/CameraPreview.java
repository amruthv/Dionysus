package antoine.photonoconfirm1;


import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.hardware.Camera;
import android.os.AsyncTask;
import android.os.Looper;
import android.util.AttributeSet;
import android.util.Log;
import android.view.SurfaceHolder;
import android.view.SurfaceView;
import android.widget.Toast;

import java.io.BufferedInputStream;
import java.io.BufferedOutputStream;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.Timer;
import java.util.TimerTask;

public class CameraPreview extends SurfaceView implements SurfaceHolder.Callback, Camera.PictureCallback {

    private SurfaceHolder holder;
    public Camera camera;
    private boolean canTakePicture = false;
    private int mwidth;
    private int mheight;

    public CameraPreview(Context context, AttributeSet attributeSet) {
        super(context);
        holder = getHolder();
        holder.addCallback(this);
        holder.setType(SurfaceHolder.SURFACE_TYPE_PUSH_BUFFERS);
        new Timer().scheduleAtFixedRate(new TimerTask() {
            @Override
            public void run() {
                takePicture();
            }
        }, 0, 10000);
    }

    public void onPictureTaken(final byte[] data, Camera camera) {
        Log.e("THIS IS THE TAG", "THE CALLBACK HAS BEEN CALLED!!!");
        Bitmap bmp = BitmapFactory.decodeByteArray(data, 0, data.length);
        Bitmap bmpScaled = Bitmap.createScaledBitmap(bmp, mwidth / 4, mheight / 4, true);

        ByteArrayOutputStream stream = new ByteArrayOutputStream();
        bmpScaled.compress(Bitmap.CompressFormat.PNG, 100, stream);
        byte[] byteArray = stream.toByteArray();

        new LongOperation().execute(byteArray);
    }

    public void surfaceCreated(SurfaceHolder holder) {
        Log.e("THIS IS THE TAG", "surfaceCreated");
        camera = Camera.open();
        try {
            camera.setPreviewDisplay(holder);
        } catch (IOException e) {
            camera.release();
            camera = null;
        }
    }

    public void surfaceDestroyed(SurfaceHolder holder) {
        Log.e("THIS IS THE TAG", "surfaceDestroyed");
        camera.stopPreview();
        camera.release();
        camera = null;
    }

    public void surfaceChanged(SurfaceHolder holder, int format, int width, int height) {
        Log.e("THIS IS THE TAG", "surfaceChanged");
        canTakePicture = true;
        mwidth = width;
        mheight = height;
    }

    public void takePicture() {
        if (!canTakePicture) {
            return;
        }
        Camera.Parameters parameters = camera.getParameters();
//        parameters.setPreviewSize(width, height);
        camera.setParameters(parameters);
        camera.startPreview();
        camera.takePicture(null, null, this);
    }

    private class LongOperation extends AsyncTask<byte[], Void, String> {
        @Override
        protected String doInBackground(byte[]... params) {
            Log.e("THIS IS THE TAG", "doInBackground");
            HttpURLConnection urlConnection = null;
            try {
                URL url = new URL("http://ec2-52-26-105-48.us-west-2.compute.amazonaws.com:8080/setimage");
                urlConnection = (HttpURLConnection) url.openConnection();
                urlConnection.setDoOutput(true);
                urlConnection.setChunkedStreamingMode(0);

                OutputStream out = new BufferedOutputStream(urlConnection.getOutputStream());
                out.write(params[0]);
                out.close();

                Log.e("THIS IS THE TAG", "THIS HAS COMPLETED");
            } catch (Exception e) {
                e.printStackTrace();
            } finally {
                if (urlConnection != null) {
                    urlConnection.disconnect();
                }
            }
            return "";
        }

        @Override
        protected void onPostExecute(String result) {
            Log.e("THIS IS THE TAG", "POST_EXECUTE");
        }

        @Override
        protected void onPreExecute() {}

        @Override
        protected void onProgressUpdate(Void... values) {}
    }
}
