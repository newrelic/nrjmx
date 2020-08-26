/*
 * Copyright 2020 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package org.newrelic.jmx;

import org.testcontainers.shaded.okhttp3.*;

public class CatsClient {

    private String baseURL;
    private OkHttpClient client = new OkHttpClient();

    public CatsClient(String baseURL) {
        this.baseURL = baseURL;
    }

    public String add(String catName) {

        RequestBody catBody = FormBody.create(
                MediaType.parse("application/json"),
                "{\"name\":\"" + catName + "\"}");

        Request request = new Request.Builder()
                .url(baseURL + "/cat")
                .post(catBody)
                .build();


        try {
            Response response = client.newCall(request).execute();
            return response.body().string();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    public String clearCats() {
        Request request = new Request.Builder()
                .url(baseURL + "/clear")
                .put(null)
                .build();

        try {
            Response response = client.newCall(request).execute();
            return response.body().string();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }
}
