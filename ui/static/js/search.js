document.body.addEventListener("htmx:configRequest", function (evt) {
    evt.detail.parameters["query"] = evt.detail.elt.value;
});
