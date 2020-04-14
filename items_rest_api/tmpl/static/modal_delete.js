// Bind click to OK button within popup
$('#confirm-delete').on('click', '.btn-ok', function(e) {

    var $modalDiv = $(e.delegateTarget);
    var id = $(this).data('recordId');
  
    $modalDiv.addClass('loading');
    $.get('/json/del?id=' + id).then(function() {
       $modalDiv.modal('hide').removeClass('loading');
       location.reload();
    });
  });
  
  // Bind to modal opening to set necessary data properties to be used to make request
  $('#confirm-delete').on('show.bs.modal', function(e) {
    var data = $(e.relatedTarget).data();
    $('.title', this).text(data.recordTitle);
    $('.btn-ok', this).data('recordId', data.recordId);
  });
