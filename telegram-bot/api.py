"""
MIT License

Copyright (c) 2021 fall2021-csc510-group40

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
"""
import requests
import uuid


class APIException(Exception):
    def __init__(self, handle, status_code, error):
        self.handle = handle
        self.status_code = status_code
        self.error = error


class APIAdaptor:
    def __init__(self, api_url, parent_logger):
        self.api_url = api_url
        self.logger = parent_logger.getChild("api")

    def get_search_results(self, title):
        self.logger.debug(f"Searching for '{title}'")

        request_data = {
            "query": {
                "title": title,
            },
        }

        self.logger.debug(f"Sending search request: {request_data}")

        rsp = requests.get(self.api_url + "/search", json=request_data, headers={"X-Request-ID": str(uuid.uuid4())})
        self.logger.debug(f"Response: {rsp.json()}")

        if not rsp.ok:
            raise APIException("search", rsp.status_code, rsp.json().get("error", None))

        return rsp.json()["results"]

    def get_format_results(self, work_id, format):
        self.logger.debug(f"Formatting {work_id} to {format}")

        rsp = requests.get(self.api_url + "/format", json={
            "work": {
                "id": work_id,
            },
            "format": format,
        }, headers={"X-Request-ID": str(uuid.uuid4())})

        if not rsp.ok:
            raise APIException("format", rsp.status_code, rsp.json().get("error", None))

        return rsp.json()["result"]
