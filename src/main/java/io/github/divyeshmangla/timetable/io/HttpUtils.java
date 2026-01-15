package io.github.divyeshmangla.timetable.io;

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
import java.security.SecureRandom;
import java.security.cert.X509Certificate;

/**
 * HTTP utilities with SSL workaround for thapar.edu certificate issues.
 */
public final class HttpUtils {
    private static final SSLContext TRUST_ALL = trustAll();

    private HttpUtils() {}

    public static InputStream download(String url) throws IOException, InterruptedException {
        var client = HttpClient.newBuilder()
                .sslContext(TRUST_ALL)
                .followRedirects(HttpClient.Redirect.NORMAL)
                .build();

        var response = client.send(
                HttpRequest.newBuilder(URI.create(url)).build(),
                HttpResponse.BodyHandlers.ofByteArray()
        );

        if (response.statusCode() != 200) {
            throw new IOException("HTTP " + response.statusCode() + " for " + url);
        }
        return new ByteArrayInputStream(response.body());
    }

    private static SSLContext trustAll() {
        try {
            var ctx = SSLContext.getInstance("TLS");
            ctx.init(null, new TrustManager[]{new TrustAllTrustManager()}, new SecureRandom());
            return ctx;
        } catch (Exception e) {
            throw new IllegalStateException("Failed to create SSL context", e);
        }
    }

    /**
     * Trust manager that accepts all certificates.
     * WARNING: This is insecure and only used as a workaround for thapar.edu certificate issues.
     */
    private static class TrustAllTrustManager implements X509TrustManager {
        @Override
        public X509Certificate[] getAcceptedIssuers() {
            return new X509Certificate[0];
        }

        @Override
        public void checkClientTrusted(X509Certificate[] chain, String authType) {
            // Trust all clients
        }

        @Override
        public void checkServerTrusted(X509Certificate[] chain, String authType) {
            // Trust all servers
        }
    }
}