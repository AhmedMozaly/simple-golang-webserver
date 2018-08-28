Vue.component("article", {
    props: ["article"],
    template: ' <div class="w3-card-4 w3-margin w3-white">\
    <img src="/Images/NMS background.jpg" alt="Nature" style="width:100%">\
    <div class="w3-container">\
      <h3><b>{{ article.Title }}</b></h3>\
      <h5>Title description, <span class="w3-opacity">{{ article.Date }}</span></h5>\
    </div>\
    <div class="w3-container">\
      <p>{{ article.PostText }}</p>\
      <div class="w3-row">\
        <div class="w3-col m8 s12">\
          <p><a href = "articles/{{ article.id }}"class="w3-button w3-padding-large w3-white w3-border"><b>READ MORE »</b></a></p>\
        </div>\
        <div class="w3-col m4 w3-hide-small">\
          <p><span class="w3-padding-large w3-right"><b>Comments  </b> <span class="w3-tag">0</span></span></p>\
        </div>\
      </div>\
    </div>\
  </div>'
  });
  //Initialize the props associated with the components
  var app = new Vue({
    el: "#app",
    data: function() {
      return {
        articles: [],
        loading: true,
        errored: false
      };
    },
    mounted () {
        axios
            .get(localhost/api/articles.json)
            .then(response => {
                this.articles = response.data
            })
            .catch(error => {
                console.log(error)
                this.errored = true;
            })
            .finally(() => this.loading = false)
    }
});
    