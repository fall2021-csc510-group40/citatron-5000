import requests


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
        # return [
        #     {
        #         "id": f"id{i}",
        #         "title": "Oasdllksa jflkjdflksa jdlkf halksdjsd lkfjaskdhflka jsdhfkalsd;flkja lsjdhf asdjlf",
        #         "authors": "Zhoka et al.",
        #     } for i in range(9)
        # ]

        request_data = {
            "query": {
                "title": title,
            },
        }

        self.logger.debug(f"Sending search request: {request_data}")

        rsp = requests.get(self.api_url + "/search", json=request_data)
        self.logger.debug(f"Response: {rsp.json()}")

        if not rsp.ok:
            raise APIException("search", rsp.status_code, rsp.json().get("error", None))

        return rsp.json()["results"]

    def get_format_results(self, work_id, format):
        self.logger.debug(f"Formatting {work_id} to {format}")
        # return f"Your {work_id} formatted to {format}"

        rsp = requests.get(self.api_url + "/format", json={
            "id": work_id,
            "format": format,
        })

        if not rsp.ok:
            raise APIException("format", rsp.status_code, rsp.json().get("error", None))

        return rsp.json()["text"]
