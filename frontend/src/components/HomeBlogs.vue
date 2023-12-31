<template>
  <v-container fluid>
    <div
      class="mt-6 mb-6 d-flex justify-space-between"
      style="font-size: 25px; margin-left: 125px"
    >
      See our blogs

      <v-dialog v-model="createDialog" width="100%">
        <template v-slot:activator="{ on, attrs }">
          <v-btn
            dark
            v-bind="attrs"
            v-on="on"
            small
            color="#252525"
            style="right: 125px"
            v-if="isLogged && !isAdmin"
          >
            <v-icon dark> mdi-plus </v-icon>
            Create
          </v-btn>
        </template>
        <v-card>
          <v-card-title class="text-h5 grey lighten-2">
            Create Blog
          </v-card-title>

          <v-card-actions>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  v-model="createForm.title"
                  label="Title"
                ></v-text-field>
              </v-col>

              <v-col cols="12">
                <v-text-field
                  v-model="createForm.avatar"
                  label="Avatar Link"
                ></v-text-field>
              </v-col>

              <v-col cols="12">
                <vue-editor v-model="createForm.content"></vue-editor>
              </v-col>
            </v-row>
          </v-card-actions>

          <v-divider></v-divider>

          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="#252525" class="white--text" @click="createItem">
              Create
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </div>

    <v-text-field
      v-model="search"
      label="Search"
      style="margin-left: 125px; margin-right: 125px; width: 200px"
      class="mb-6"
      append-icon="mdi-magnify"
    ></v-text-field>

    <div style="margin-left: 125px; margin-right: 125px" class="mb-6">
      <v-row>
        <v-col cols="3" class="post" v-for="(item, i) in items" :key="i">
          <div class="post-avatar" @click="viewBlog(item.id)">
            <v-img
              width="100%"
              height="200"
              class="rounded-lg"
              :src="item.avatar"
            ></v-img>
          </div>
          <div class="post-title" @click="viewBlog(item.id)">
            <strong>
              {{ item.title }}
            </strong>
          </div>

          <div>
            <p style="opacity: 0.7; font-size: 12px">{{ item.created_at }}</p>
          </div>
        </v-col>
      </v-row>

      <div class="more-post">
        <v-icon
          dark
          v-if="!loading"
          color="#252525"
          class="icon"
          @click="morePost"
        >
          mdi-chevron-down</v-icon
        >

        <v-progress-circular
          v-if="loading"
          indeterminate
          color="primary"
        ></v-progress-circular>
      </div>
    </div>
  </v-container>
</template>

<script>
import { VueEditor } from "vue2-editor";
export default {
  props: {
    userInfor: Object,
    isLogged: Boolean,
    isAdmin: Boolean,
  },

  components: {
    VueEditor,
  },

  data() {
    return {
      search: "",
      itemsPerPage: 4,
      items: [],
      loading: false,
      createForm: {
        title: "",
        avatar: "",
        content: "",
      },

      createDialog: false,
    };
  },

  watch: {
    search(newValue) {
      this.loadItems();
    },
  },

  methods: {
    loadItems() {
      this.loading = true;
      setTimeout(() => {
        this.loading = false;
        this.$axios
          .get(`/blogs?page=1&limit=${this.itemsPerPage}&search=${this.search}`)
          .then((res) => {
            this.items = res.data.items || [];
          })
          .catch((err) => {
            console.error(err);
          });
      }, 1000);
    },

    viewBlog(id) {
      this.$router.push({ name: "blog-page", params: { id: id } });
    },

    morePost() {
      this.itemsPerPage += 4;
      this.loadItems();
    },

    createItem() {
      this.$axios
        .post(`/blogs`, {
          title: this.createForm.title,
          avatar: this.createForm.avatar,
          content: this.createForm.content,
          user_id: this.userInfor.id,
        })

        .then((res) => {
          this.createDialog = false;
          this.createForm = {
            title: "",
            avatar: "",
            content: "",
          };

          this.loadItems();
        })

        .catch((err) => {
          console.error(err);
        });
    },
  },

  created() {
    this.loadItems();
  },
};
</script>

<style lang="scss" scoped>
.post {
  .post-avatar,
  .post-title {
    cursor: pointer;
    &:hover {
      opacity: 0.5;
    }
  }
}

.more-post {
  text-align: center;

  .icon {
    cursor: pointer;
    &:hover {
      opacity: 0.5;
    }
  }
}
</style>
