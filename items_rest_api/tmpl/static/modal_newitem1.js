$(document).ready(function(){	
	$("#contactForm").submit(function(event){
		submitForm();
		return false;
	});
});


function submitForm(){

        var jsonData = {};
        $.each($('#contactForm').serializeArray(), function() {
            jsonData[this.name] = this.value;
         });

        var datajson = JSON.stringify(jsonData);

        $.ajax({
            type: "POST",
            url: "/add",
            contentType: 'application/json;charset=UTF-8',
            dataType: 'json',
            data: datajson,
            success: function(response){
                $("#contact").html(response)
                $("#contact-modal").modal('hide');
                location.reload();
            },
            error: function(){
                alert("Error");
            }
        }); 
}
