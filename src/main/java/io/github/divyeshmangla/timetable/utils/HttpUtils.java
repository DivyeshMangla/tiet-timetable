package io.github.divyeshmangla.timetable.utils;

import javax.net.ssl.SSLContext;
import javax.net.ssl.TrustManager;
import javax.net.ssl.X509TrustManager;
import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.security.KeyManagementException;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;
import java.security.cert.X509Certificate;
import java.time.Duration;

/**
 * HTTP utilities with custom SSL handling for thapar.edu certificate issues.
 */
public final class HttpUtils {
    private static final String TRUSTED_DOMAIN = "https://www.thapar.edu/upload/files/";
    private static final SSLContext TRUST_ALL_CONTEXT = createTrustAllContext();

    private HttpUtils() {}

    public static InputStream download(String url) throws IOException, InterruptedException {
        validateTrustedUrl(url);

        HttpClient client = HttpClient
                .newBuilder()
                .connectTimeout(Duration.ofSeconds(20))
                .followRedirects(HttpClient.Redirect.NORMAL)
                .sslContext(TRUST_ALL_CONTEXT)
                .build();

        HttpRequest request = HttpRequest
                .newBuilder()
                .uri(URI.create(url))
                .GET()
                .build();

        HttpResponse<byte[]> response = client.send(request, HttpResponse.BodyHandlers.ofByteArray());

        if (response.statusCode() != 200) {
            throw new IOException("HTTP " + response.statusCode() + " for " + url);
        }

        return new ByteArrayInputStream(response.body());
    }

    /**
     * Creates an SSLContext that trusts all certificates (workaround for thapar.edu)
     */
    private static SSLContext createTrustAllContext() {
        try {
            TrustManager[] trustAllCerts = new TrustManager[]{
                    new X509TrustManager() {
                        public X509Certificate[] getAcceptedIssuers() {
                            return new X509Certificate[0];
                        }
                        public void checkClientTrusted(X509Certificate[] certs, String authType) {} //NOSONAR
                        public void checkServerTrusted(X509Certificate[] certs, String authType) {} //NOSONAR
                    }
            };

            SSLContext sslContext = SSLContext.getInstance("TLS");
            sslContext.init(null, trustAllCerts, new SecureRandom());
            return sslContext;
        } catch (NoSuchAlgorithmException | KeyManagementException e) {
            throw new IllegalStateException("Failed to create SSL context", e);
        }
    }

    private static void validateTrustedUrl(String url) {
        if (!url.startsWith(TRUSTED_DOMAIN)) {
            throw new SecurityException("URL must start with " + TRUSTED_DOMAIN + " but was: " + url);
        }
    }
}