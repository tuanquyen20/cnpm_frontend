<template>
  <v-container fluid>
    <div>
      <v-dialog v-model="createDialog" width="100%">
        <template v-slot:activator="{ on, attrs }">
          <v-btn
            dark
            v-bind="attrs"
            v-on="on"
            small
            color="#252525"
            class="mb-4"
          >
            <v-icon dark> mdi-plus </v-icon>
            Create
          </v-btn>
        </template>
        <v-card>
          <v-card-title class="text-h5 grey lighten-2">
            Create Blogs
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

      <v-text-field
        v-model="search"
        label="Search"
        append-icon="mdi-magnify"
      ></v-text-field>

      <v-data-table
        :headers="headers"
        :items="items"
        :page.sync="page"
        :items-per-page="itemsPerPage"
        hide-default-footer
        class="elevation-1"
        :loading="loading"
      >
        <template v-slot:[`item.title`]="{ item }">
          <strong>{{ item.title }}</strong>
        </template>

        <template v-slot:[`item.actions`]="{ item }">
          <v-icon small class="mr-2" @click="showUpdateForm(item)">
            mdi-pencil
          </v-icon>

          <v-icon small @click="deleteItem(item)"> mdi-delete </v-icon>
        </template>
      </v-data-table>

      <v-pagination
        v-model="page"
        :length="pageCount"
        @next="loadItems"
        @previous="loadItems"
        @input="loadItems"
      ></v-pagination>
    </div>

    <v-dialog v-model="editDialog" width="100%">
      <v-card>
        <v-card-title class="text-h5 grey lighten-2">
          Update Post
        </v-card-title>

        <v-card-actions>
          <v-row>
            <v-col cols="12">
              <v-text-field
                v-model="updateForm.title"
                label="Title"
              ></v-text-field>
            </v-col>

            <v-col cols="12">
              <v-text-field
                v-model="updateForm.avatar"
                label="Avatar Link"
              ></v-text-field>
            </v-col>

            <v-col cols="12">
              <vue-editor v-model="updateForm.content"></vue-editor>
            </v-col>
          </v-row>
        </v-card-actions>

        <v-divider></v-divider>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="#252525" class="white--text" @click="updateItem">
            Update
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script>
import { VueEditor } from "vue2-editor";

export default {
  components: {
    VueEditor,
  },

  data() {
    return {
      search: "",
      userInfor: {},
      isLogged: false,
      isAdmin: false,
      editDialog: false,
      createDialog: false,
      loading: false,
      page: 1,
      pageCount: null,
      itemsPerPage: 10,
      items: [],
      headers: [
        { text: "ID", value: "id" },
        { text: "Title", value: "title" },
        { text: "Created At", value: "created_at" },
        { text: "Actions", value: "actions", sortable: false },
      ],

      updateForm: {
        id: null,
        created_at: "",
        title: "",
        avatar: "",
        content: "",
      },

      createForm: {
        title: "",
        avatar: "",
        content: "",
      },
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
          .get(
            `/blogs?page=${this.page}&limit=${this.itemsPerPage}&search=${this.search}`
          )
          .then((res) => {
            this.items = res.data.items || [];
            this.pageCount = res.data.page_count;
          })
          .catch((err) => {
            console.error(err);
          });
      }, 1000);
    },

    deleteItem(item) {
      this.$axios
        .delete(`/blogs/${item.id}`)
        .then((res) => {
          this.loadItems();
        })
        .catch((err) => {
          console.error(err);
        });
    },

    showUpdateForm(item) {
      this.editDialog = true;
      this.updateForm = {
        id: item.id,
        created_at: item.created_at,
        title: item.title,
        avatar: item.avatar,
        content: item.content,
        user_id: item.user_id,
      };
    },

    updateItem() {
      this.$axios
        .put(`/blogs/${this.updateForm.id}`, {
          title: this.updateForm.title,
          avatar: this.updateForm.avatar,
          content: this.updateForm.content,
          created_at: this.updateForm.created_at,
          user_id: this.updateForm.user_id,
        })

        .then((res) => {
          this.editDialog = false;
          this.loadItems();
        })
        .catch((err) => {
          console.error(err);
        });
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
    this.userInfor = this.$store.state.userInfor;
    this.isAdmin = this.$store.state.isAdmin;
    this.isLogged = this.$store.state.isLogged;
    this.loadItems();
  },
};
</script>
