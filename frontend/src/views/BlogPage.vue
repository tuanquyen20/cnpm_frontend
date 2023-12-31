<template>
  <v-container fluid>
    <div class="blog">
      <div>
        <v-icon @click="backToHome" color="blue" class="blog-back-to-home"
          >mdi-keyboard-backspace</v-icon
        >
        <span @click="backToHome" class="ml-2 blue--text blog-back-to-home"
          >Back to home</span
        >
      </div>

      <div class="author mt-6 mb-6 d-flex justify-center align-center">
        <v-img
          v-if="authorInfor.avatar_data"
          class="rounded-pill mr-4"
          max-height="200px"
          max-width="200px"
          :src="'data:image/jpeg;base64,' + authorInfor.avatar_data"
        ></v-img>

        <div class="infor ml-4">
          <strong>{{ authorInfor.name }}</strong>
          <p>{{ authorInfor.email }}</p>
        </div>
      </div>

      <p class="blog-title">{{ item.title }}</p>
      <p class="blog-time">{{ item.created_at }}</p>
      <v-divider></v-divider>

      <div v-html="item.content" class="mt-6"></div>
    </div>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      item: {},
      authorInfor: {},
    };
  },

  methods: {
    async getBlog(id) {
      await this.$axios
        .get(`/blogs/${id}`)
        .then((res) => {
          this.item = res.data;
        })
        .catch((err) => {
          console.error(err);
        });

      this.getAuthor();
    },

    getAuthor() {
      this.$axios
        .get(`/accounts/${this.item.user_id}`)
        .then((res) => {
          this.authorInfor = res.data;
        })
        .catch((err) => {
          console.error(err);
        });
    },

    backToHome() {
      this.$router.push({ name: "home" });
    },
  },

  created() {
    this.getBlog(this.$route.params.id);
  },
};
</script>

<style lang="scss">
img {
    max-width: 100%;
    border-radius: 5px;
  }

  iframe {
    width: 100%;
    min-height: 300px;
  }


.blog {
  width: 750px;
  margin: 100px auto auto auto;

  img {
    max-width: 100%;
    border-radius: 5px;
  }

  iframe {
    width: 100%;
    min-height: 300px;
  }

  .blog-back-to-home {
    cursor: pointer;
  }

  .blog-title {
    font-size: 30px;
    font-weight: bold;
  }

  .blog-time {
    font-size: 13px;
    opacity: 0.5;
  }
}
</style>
