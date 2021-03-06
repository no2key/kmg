package kmgRpcJava

import (
	"bytes"
)

func tplGenerateCode(config *tplConfig) string {
	var _buf bytes.Buffer
	_buf.WriteString(`
package `)
	_buf.WriteString(config.OutPackageName)
	_buf.WriteString(`;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.google.gson.JsonDeserializationContext;
import com.google.gson.JsonDeserializer;
import com.google.gson.JsonElement;
import com.google.gson.JsonParseException;
import com.google.gson.JsonPrimitive;
import com.google.gson.JsonSerializationContext;
import com.google.gson.JsonSerializer;
import com.google.gson.JsonSyntaxException;

import javax.crypto.Cipher;
import javax.crypto.spec.IvParameterSpec;
import javax.crypto.spec.SecretKeySpec;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.lang.reflect.Type;
import java.net.HttpURLConnection;
import java.net.URL;
import java.nio.charset.Charset;
import java.security.MessageDigest;
import java.security.SecureRandom;
import java.util.Arrays;
import java.util.Calendar;
import java.util.Date;
import java.util.List;
import java.util.SimpleTimeZone;
import java.util.TimeZone;
import java.util.zip.Deflater;
import java.util.zip.Inflater;

/*
    example:
        RpcDemo.ConfigDefaultClient("http://127.0.0.1:34567","abc psk") // added in some init function.
        String result = RpcDemo.GetDefaultClient().PostScoreInt("abc",1) // use the rpc everywhere.
 */

public class `)
	_buf.WriteString(config.ClassName)
	_buf.WriteString(` {
    //类型列表
    `)
	for _, innerClass := range config.InnerClassList {
		_buf.WriteString(`
        `)
		_buf.WriteString(innerClass.tplInnerClass())
		_buf.WriteString(`
    `)
	}
	_buf.WriteString(`

    public static class Client{
        // 所有Api列表
        `)
	for _, api := range config.ApiList {
		_buf.WriteString(`
            `)
		_buf.WriteString(api.tplApiClient())
		_buf.WriteString(`
        `)
	}
	_buf.WriteString(`


        //引入的不会变的库代码.还需要com.google.gson 这个package的依赖
        public String RemoteUrl;
        public byte[] Psk;
        private <T> T sendRequest(String apiName,Object reqData,Class<T> tClass) throws Exception{
            String inDataString = kmgJson.MarshalToString(reqData); // UTF8? 啥编码?
            if (apiName.length()>255){
                throw new Exception("len(apiName)>255");
            }
            ByteArrayOutputStream baos = new ByteArrayOutputStream();
            baos.write(apiName.length());
            baos.write(KmgString.StringToArrayByte(apiName));
            baos.write(KmgString.StringToArrayByte(inDataString));
            byte[] inByte = baos.toByteArray();
            if (this.Psk!=null){
                inByte = kmgCrypto.CompressAndEncryptBytesEncodeV2(this.Psk, inByte);
            }
            byte[] outBytes =  kmgHttp.SimplePost(this.RemoteUrl, inByte);
            outBytes = kmgCrypto.CompressAndEncryptBytesDecodeV2(this.Psk, outBytes);
            if (outBytes.length==0){
                throw new Exception("outBytes.length==0");
            }
            String AfterString = KmgString.ArrayByteToString(Arrays.copyOfRange(outBytes,1,outBytes.length));
            if (outBytes[0]==1){ //error
                throw new Exception(AfterString);
            }else if (outBytes[0]==2) { //success
                return kmgJson.UnmarshalFromString(AfterString, tClass);
            }
            throw new Exception("httpjsonApi protocol error 1 "+outBytes[0]);
        }
    }
    public static void ConfigDefaultClient(String RemoteUrl,String pskStr){
        defaultClient = new Client();
        defaultClient.RemoteUrl = RemoteUrl;
        defaultClient.Psk = kmgCrypto.Get32PskFromString(pskStr);
    }
    private static Client defaultClient;
    public static Client GetDefaultClient(){
        return defaultClient;
    }
    public static class KmgString{
        public static final Charset UTF_8 = Charset.forName("UTF-8");
        public static byte[] StringToArrayByte(String str){
            return str.getBytes(UTF_8);
        }
        public static String ArrayByteToString(byte[] bytes){
            return new String(bytes, UTF_8);
        }
    }
    public static class kmgBytes {
        public static byte[] Slice(byte[] in,int start,int end){
            return Arrays.copyOfRange(in,start,end);
        }
    }
    public static class kmgIo {
        private static final int EOF = -1;
        public static byte[] InputStreamReadAll(final InputStream input) throws IOException {
            final ByteArrayOutputStream output = new ByteArrayOutputStream();
            copy(input, output);
            return output.toByteArray();
        }
        public static int copy(final InputStream input, final OutputStream output) throws IOException {
            final long count = copyLarge(input, output,new byte[8192]);
            if (count > Integer.MAX_VALUE) {
                return -1;
            }
            return (int) count;
        }
        public static long copyLarge(final InputStream input, final OutputStream output, final byte[] buffer)
                throws IOException {
            long count = 0;
            int n = 0;
            while (EOF != (n = input.read(buffer))) {
                output.write(buffer, 0, n);
                count += n;
            }
            return count;
        }
    }
    public static class kmgHttp {
        public static byte[] SimplePost(String urls,byte[] inByte) throws Exception{
            System.setProperty("http.keepAlive", "false");
            URL url = new URL(urls);
            HttpURLConnection conn = (HttpURLConnection)url.openConnection();
            conn.setRequestMethod("POST");
            conn.setRequestProperty("Content-Type", "image/jpeg");
            conn.setUseCaches(false);
            conn.setDoInput(true);
            conn.setDoOutput(true);
            conn.setReadTimeout(10000);
            conn.setConnectTimeout(5000);
            OutputStream os = conn.getOutputStream();
            os.write(inByte);
            os.flush();
            os.close();
            InputStream is;
            if (conn.getResponseCode()==200){
                is = conn.getInputStream();
            }else{
                is = conn.getErrorStream();
            }
            byte[] outByte = kmgIo.InputStreamReadAll(is);
            is.close();
            return outByte;
        }
    }
    public static class kmgCrypto {
        private static byte[] magicCode4 = new byte[]{(byte)0xa7,(byte)0x97,0x6d,0x15};
        // key lenth 32
        public static byte[] CompressAndEncryptBytesEncodeV2(byte[] key,byte[] data) throws Exception{
            data = compressV2(data);
            byte[] cbcIv = KmgRand.MustCryptoRandBytes(16);
            ByteArrayOutputStream buf = new ByteArrayOutputStream();
            buf.write(data);
            buf.write(magicCode4);
            Cipher cipher = Cipher.getInstance("AES/CTR/NoPadding");
            cipher.init(Cipher.ENCRYPT_MODE, new SecretKeySpec(key, "AES"), new IvParameterSpec(cbcIv));
            byte[] encrypted = cipher.doFinal(buf.toByteArray());
            buf = new ByteArrayOutputStream();
            buf.write(cbcIv);
            buf.write(encrypted);
            return buf.toByteArray();
        }
        // key lenth 32
        public static byte[] CompressAndEncryptBytesDecodeV2(byte[] key,byte[] data) throws Exception {
            if (data.length < 21) {
                throw new Exception("[kmgCrypto.CompressAndEncryptBytesDecode] input data too small");
            }
            byte[] cbcIv = kmgBytes.Slice(data, 0, 16);
            byte[] encrypted = kmgBytes.Slice(data, 16, data.length);
            Cipher cipher = Cipher.getInstance("AES/CTR/NoPadding");
            cipher.init(Cipher.DECRYPT_MODE, new SecretKeySpec(key, "AES"), new IvParameterSpec(cbcIv));
            byte[] decrypted = cipher.doFinal(encrypted);
            byte[] compressed = kmgBytes.Slice(decrypted, 0, decrypted.length - 4);
            if (!Arrays.equals(magicCode4, kmgBytes.Slice(decrypted, decrypted.length - 4, decrypted.length))){
                throw new Exception("[kmgCrypto.CompressAndEncryptBytesDecode] magicCode not match");
            }
            return uncompressV2(compressed);
        }
        private static byte[] compressV2(byte[] data) throws Exception{
            byte[] outData = kmgCompress.ZlibMustCompress(data);
            if (outData.length>=data.length){
                ByteArrayOutputStream buf = new ByteArrayOutputStream();
                buf.write(0);
                buf.write(data);
                return buf.toByteArray();
            }else{
                ByteArrayOutputStream buf = new ByteArrayOutputStream();
                buf.write(1);
                buf.write(outData);
                return buf.toByteArray();
            }
        }
        private static byte[] uncompressV2(byte[] data) throws Exception{
            if (data.length==0){
                throw new Exception("[uncopressV2] len(inData)==0");
            }
            if (data[0]==0){
                return kmgBytes.Slice(data, 1, data.length);
            }
            return kmgCompress.ZlibUnCompress(kmgBytes.Slice(data, 1, data.length));
        }
        public static byte[] Sha512Sum(byte[] data) {
            try {
                MessageDigest sh = MessageDigest.getInstance("SHA-512");
                sh.update(data);
                return sh.digest();
            }catch(Exception e){
                System.out.println(e.getMessage());
                e.printStackTrace();
            }
            return null;
        }
        public static byte[] Get32PskFromString(String s){
            return kmgBytes.Slice(kmgCrypto.Sha512Sum(KmgString.StringToArrayByte(s)), 0, 32);
        }
    }
    public static class kmgCompress {
        public static byte[] ZlibMustCompress(byte[] inB){
            Deflater deflater = new Deflater();
            deflater.setInput(inB);
            deflater.finish();
            ByteArrayOutputStream baos = new ByteArrayOutputStream();
            byte[] buf = new byte[8192];
            while (!deflater.finished()) {
                int byteCount = deflater.deflate(buf);
                baos.write(buf, 0, byteCount);
            }
            deflater.end();
            byte[] out = baos.toByteArray();
            return out;
        }
        public static byte[] ZlibUnCompress(byte[] inB) throws Exception{
            Inflater deflater = new Inflater();
            deflater.setInput(inB);
            ByteArrayOutputStream baos = new ByteArrayOutputStream();
            byte[] buf = new byte[8192];
            while (!deflater.finished()) {
                int byteCount = deflater.inflate(buf);
                if (byteCount==0){
                    break;
                }
                baos.write(buf, 0, byteCount);
            }
            deflater.end();
            return baos.toByteArray();
        }
    }
    public static class kmgSync {
        public static class Once{
            private Object locker = new Object();
            private boolean isInit = false;
            public void Do(Runnable f){
                synchronized (locker){
                    if (isInit){
                        return;
                    }
                    f.run();
                    isInit = true;
                }
            }
        }
    }
    public static class kmgJson {
        public static String MarshalToString(Object data){
            return getGson().toJson(data);
        }
        public static<T> T UnmarshalFromString(String s,Class<T> t) throws JsonSyntaxException {
            if (t==void.class){
                return null;
            }
            return getGson().fromJson(s,t);
        }
        private static Gson gson;
        private static kmgSync.Once gsonOnce = new kmgSync.Once();
        private static Gson getGson(){
            gsonOnce.Do(new Runnable() {
                @Override
                public void run() {
                    JsonSerializer<Date> ser = new JsonSerializer<Date>() {
                        @Override
                        public JsonElement serialize(Date src, Type typeOfSrc, JsonSerializationContext
                                context) {
                            if (src == null) {
                                return null;
                            } else {
                                return new JsonPrimitive(KmgTime.FormatGolangDate(src));
                            }
                        }
                    };
                    JsonDeserializer<Date> deser = new JsonDeserializer<Date>() {
                        @Override
                        public Date deserialize(JsonElement json, Type typeOfT,
                                                JsonDeserializationContext context) throws JsonParseException {
                            if (json == null) {
                                return null;
                            } else {
                                try {
                                    return KmgTime.ParseGolangDate(json.getAsString());
                                } catch (Exception e) {
                                    throw new JsonParseException(e.getMessage());
                                }
                            }
                        }
                    };
                    gson = new GsonBuilder()
                            .registerTypeAdapter(Date.class, ser)
                            .registerTypeAdapter(Date.class, deser).create();
                }
            });
            return gson;
        }
    }
    public static class KmgRand{
        public static byte[] MustCryptoRandBytes(int length) {
            SecureRandom sr = new SecureRandom();
            byte[] output = new byte[length];
            sr.nextBytes(output);
            return output;
        }
    }
    public static class KmgTime{
        public static Date ParseGolangDate(String st) throws Exception{
            Calendar cal = Calendar.getInstance();
            int year = Integer.parseInt(st.substring(0, 4));
            int month = Integer.parseInt(st.substring(5,7));
            int day = Integer.parseInt(st.substring(8,10));
            int hour = Integer.parseInt(st.substring(11, 13));
            int minute = Integer.parseInt(st.substring(14,16));
            int second = Integer.parseInt(st.substring(17,19));
            float nonaSecond = 0;
            int tzStartPos = 19;
            if (st.charAt(19)=='.'){
                // 从19开始找到第一个不是数字的字符串.
                for (;;){
                    tzStartPos++;
                    if (st.length()<=tzStartPos){
                        throw new Exception("can not parse "+st);
                    }
                    char thisChar = st.charAt(tzStartPos);
                    if (thisChar>='0' && thisChar<='9'){
                    }else{
                        break;
                    }
                }
                nonaSecond = Float.parseFloat("0." + st.substring(20, tzStartPos));
            }
            cal.set(Calendar.MILLISECOND,(int)(nonaSecond*1e3));
            char tzStart = st.charAt(tzStartPos);
            if (tzStart=='Z'){
                cal.setTimeZone(TimeZone.getTimeZone("UTC"));
            }else {
                int tzHour = Integer.parseInt(st.substring(tzStartPos+1,tzStartPos+3));
                int tzMin = Integer.parseInt(st.substring(tzStartPos+4,tzStartPos+6));
                int tzOffset = tzHour*3600*1000 + tzMin * 60*1000;
                if (tzStart=='-'){
                    tzOffset = - tzOffset;
                }
                TimeZone tz = new SimpleTimeZone(tzOffset,"");
                cal.setTimeZone(tz);
            }
            cal.set(year,month-1,day,hour,minute,second);
            return cal.getTime();
        }
        public static String FormatGolangDate(Date date){
            Calendar cal = Calendar.getInstance();
            cal.setTime(date);
            StringBuilder buf = new StringBuilder();
            formatYear(cal,buf);
            buf.append('-');
            formatMonth(cal, buf);
            buf.append('-');
            formatDays(cal, buf);
            buf.append('T');
            formatHours(cal, buf);
            buf.append(':');
            formatMinutes(cal, buf);
            buf.append(':');
            formatSeconds(cal,buf);
            formatTimeZone(cal,buf);
            return buf.toString();
        }
        private static void formatYear(Calendar cal, StringBuilder buf) {
            int year = cal.get(Calendar.YEAR);
            String s;
            if (year <= 0) // negative value
            {
                s = Integer.toString(1 - year);
            } else // positive value
            {
                s = Integer.toString(year);
            }
            while (s.length() < 4) {
                s = '0' + s;
            }
            if (year <= 0) {
                s = '-' + s;
            }
            buf.append(s);
        }
        private static void formatMonth(Calendar cal, StringBuilder buf) {
            formatTwoDigits(cal.get(Calendar.MONTH) + 1, buf);
        }
        private static void formatDays(Calendar cal, StringBuilder buf) {
            formatTwoDigits(cal.get(Calendar.DAY_OF_MONTH), buf);
        }
        private static void formatHours(Calendar cal, StringBuilder buf) {
            formatTwoDigits(cal.get(Calendar.HOUR_OF_DAY), buf);
        }
        private static void formatMinutes(Calendar cal, StringBuilder buf) {
            formatTwoDigits(cal.get(Calendar.MINUTE), buf);
        }
        private static void formatSeconds(Calendar cal, StringBuilder buf) {
            formatTwoDigits(cal.get(Calendar.SECOND), buf);
            if (cal.isSet(Calendar.MILLISECOND)) { // milliseconds
                int n = cal.get(Calendar.MILLISECOND);
                if (n != 0) {
                    String ms = Integer.toString(n);
                    while (ms.length() < 3) {
                        ms = '0' + ms; // left 0 paddings.
                    }
                    buf.append('.');
                    buf.append(ms);
                }
            }
        }
        /** formats time zone specifier. */
        private static void formatTimeZone(Calendar cal, StringBuilder buf) {
            TimeZone tz = cal.getTimeZone();
            if (tz == null) {
                return;
            }
            // otherwise print out normally.
            int offset = tz.getOffset(cal.getTime().getTime());
            if (offset == 0) {
                buf.append('Z');
                return;
            }
            if (offset >= 0) {
                buf.append('+');
            } else {
                buf.append('-');
                offset *= -1;
            }
            offset /= 60 * 1000; // offset is in milli-seconds
            formatTwoDigits(offset / 60, buf);
            buf.append(':');
            formatTwoDigits(offset % 60, buf);
        }
        /** formats Integer into two-character-wide string. */
        private static void formatTwoDigits(int n, StringBuilder buf) {
            // n is always non-negative.
            if (n < 10) {
                buf.append('0');
            }
            buf.append(n);
        }
    }
}


`)
	return _buf.String()
}
