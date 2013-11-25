function executeCode(i) {
    var slideNumber = "#kuvCode" + i;
    var sourceCode = $(slideNumber).eq(0).val();
    var parameters = {
        "a": JSON.stringify({
            Code: sourceCode
        })
    };

    $.getJSON('http://localhost:8080', parameters)
        .done(function(str) {
            document.getElementById('kuvOutput' + i).innerHTML = str.Output;
        })
        .fail(function(e) {
            alert("Error: Check if the compiler server is running correctly.");
        });
}
