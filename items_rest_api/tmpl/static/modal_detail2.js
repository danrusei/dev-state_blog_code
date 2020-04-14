$('#confirm-detail').on('show.bs.modal', function(e) {
    var data = $(e.relatedTarget).data();
    $('.id', this).text(data.recordId);
    $('.created', this).text(data.recordCreated);
    $('.title', this).text(data.recordTitle);
    $('.expdate', this).text(data.recordExpdate);
    $('.expopen', this).text(data.recordExpopen);
    $('.targetage', this).text(data.recordTargetage);
    $('.comment', this).text(data.recordComment);
    $('.isopen', this).text(data.recordIsopen);
    $('.opened', this).text(data.recordOpened);
    $('.isvalid', this).text(data.recordIsvalid);
    $('.daysvalid', this).text(data.recordDaysvalid);
  });
