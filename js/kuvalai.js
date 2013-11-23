function executeCode(i) {
	var slideNumber = "#kuvCode"+i;
	var sourceCode = $(slideNumber).eq(0).val();
	alert(sourceCode);
	var parameters = {"a" : JSON.stringify({ Code: sourceCode})}; 
	alert(parameters);

	$.getJSON('http://localhost:8080', parameters)
		.done(function(str){
			document.getElementById('kuvOutput1').innerHTML = str.Output;
		})
	.fail(function(e) {
		alert("Failure");
	});
}
