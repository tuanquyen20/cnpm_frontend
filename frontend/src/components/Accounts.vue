<template>
  <v-container fluid>
    <div>
      <v-dialog v-model="createDialog" width="700">
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
            Create User
          </v-card-title>

          <v-card-actions>
            <v-row>
              <v-col cols="6">
                <v-img
                  v-if="createForm.avatarData"
                  class="rounded-pill"
                  width="200"
                  height="200"
                  :src="createForm.avatarData"
                ></v-img>
              </v-col>

              <v-col cols="6">
                <v-file-input
                  v-model="createForm.avatarFile"
                  accept="image/*"
                  label="Avatar file"
                  append-icon="mdi-camera"
                  @change="changeAvatarFile('create')"
                  :clearable="false"
                ></v-file-input>
              </v-col>

              <v-col cols="12">
                <v-text-field
                  v-model="createForm.name"
                  label="Name"
                ></v-text-field>
              </v-col>

              <v-col cols="12">
                <v-text-field
                  v-model="createForm.phone_number"
                  label="Phone Number"
                ></v-text-field>
              </v-col>

              <v-col cols="12">
                <v-text-field
                  v-model="createForm.email"
                  label="Email"
                ></v-text-field>
              </v-col>

              <v-col cols="12">
                <v-text-field
                  v-model="createForm.password"
                  label="Password"
                ></v-text-field>
              </v-col>

              <v-col cols="12">
                <v-select
                  v-model="createForm.type"
                  :items="['admin', 'customer']"
                  label="Type"
                  single-line
                ></v-select>
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

    <v-dialog v-model="editDialog" width="700">
      <v-card>
        <v-card-title class="text-h5 grey lighten-2">
          Update account
        </v-card-title>

        <v-card-actions>
          <v-row>
            <v-col cols="6">
              <v-img
                v-if="!updateForm.avatarDataTemp"
                class="rounded-pill"
                width="200"
                height="200"
                :src="'data:image/jpeg;base64,' + updateForm.avatarData"
              ></v-img>

              <v-img
                v-if="updateForm.avatarDataTemp"
                class="rounded-pill"
                width="200"
                height="200"
                :src="updateForm.avatarDataTemp"
              ></v-img>
            </v-col>

            <v-col cols="6">
              <v-file-input
                v-model="updateForm.avatarFile"
                @change="changeAvatarFile('update')"
                :clearable="false"
                accept="image/*"
                label="Avatar file"
                prepend-icon="mdi-camera"
              ></v-file-input>
            </v-col>

            <v-col cols="12">
              <v-text-field
                v-model="updateForm.name"
                label="Name"
              ></v-text-field>
            </v-col>

            <v-col cols="12">
              <v-text-field
                v-model="updateForm.phone_number"
                label="Phone Number"
              ></v-text-field>
            </v-col>

            <v-col cols="12">
              <v-text-field
                v-model="updateForm.email"
                label="Email"
              ></v-text-field>
            </v-col>

            <v-col cols="12">
              <v-text-field
                v-model="updateForm.password"
                label="Password"
              ></v-text-field>
            </v-col>

            <v-col cols="12">
              <v-select
                v-model="updateForm.type"
                :items="['admin', 'customer']"
                label="Type"
                single-line
              ></v-select>
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
export default {
  data() {
    return {
      search: "",

      editDialog: false,
      createDialog: false,
      loading: false,
      page: 1,
      pageCount: 0,
      itemsPerPage: 10,
      items: [],
      headers: [
        { text: "ID", value: "id" },
        { text: "Email", value: "email" },
        { text: "Name", value: "name" },
        { text: "Phone Number", value: "phone_number" },
        { text: "Type", value: "type" },
        { text: "Actions", value: "actions", sortable: false },
      ],

      updateForm: {
        id: null,
        name: "",
        email: "",
        password: "",
        phone_number: "",
        avatarData: null,
        avatarName: null,
        avatarFile: null,
        avatarDataTemp: null,
      },

      createForm: {
        name: "",
        email: "",
        password: "",
        phone_number: "",
        avatarData: null,
        avatarFile: null,
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
            `/accounts?page=${this.page}&limit=${this.itemsPerPage}&search=${this.search}`
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
        .delete(`/accounts/${item.id}`)
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
        name: item.name,
        email: item.email,
        password: item.password,
        phone_number: item.phone_number,
        type: item.type,
        avatarName: item.avatar_name,
        avatarData: item.avatar_data,
        avatarFile: null,
        avatarDataTemp: null,
      };
    },

    updateItem() {
      if (this.updateForm.avatarFile) {
        const formData = new FormData();
        formData.append("image", this.updateForm.avatarFile);
        formData.append(
          "account",
          JSON.stringify({
            name: this.updateForm.name,
            phone_number: this.updateForm.phone_number,
            email: this.updateForm.email,
            password: this.updateForm.password,
            type: this.updateForm.type,
          })
        );

        this.$axios
          .put(`/accounts/${this.updateForm.id}`, formData)
          .then((res) => {
            this.editDialog = false;
            this.updateForm = {
              id: null,
              name: "",
              email: "",
              password: "",
              phone_number: "",
              avatarData: null,
              avatarName: null,
              avatarFile: null,
              avatarDataTemp: null,
            };

            this.loadItems();
          })
          .catch((err) => {
            console.error(err);
          });
      } else {
        this.$axios
          .put(`/accounts-without-avatar/${this.updateForm.id}`, {
            name: this.updateForm.name,
            phone_number: this.updateForm.phone_number,
            email: this.updateForm.email,
            password: this.updateForm.password,
            type: this.updateForm.type,
          })
          .then((res) => {
            this.editDialog = false;
            this.updateForm = {
              id: null,
              name: "",
              email: "",
              password: "",
              phone_number: "",
              avatarData: null,
              avatarName: null,
              avatarFile: null,
              avatarDataTemp: null,
            };

            this.loadItems();
          })
          .catch((err) => {
            console.error(err);
          });
      }
    },

    changeAvatarFile(type) {
      if (type === "create") {
        if (this.createForm.avatarFile) {
          this.readImageFile(this.createForm.avatarFile, "create");
        }
      } else {
        if (this.updateForm.avatarFile) {
          this.readImageFile(this.updateForm.avatarFile, "update");
        }
      }
    },

    readImageFile(file, type) {
      const reader = new FileReader();
      reader.onload = (e) => {
        if (type === "create") {
          this.createForm.avatarData = e.target.result;
        } else {
          this.updateForm.avatarDataTemp = e.target.result;
        }
      };
      reader.readAsDataURL(file);
    },

    createItem() {
      const formData = new FormData();
      formData.append("image", this.createForm.avatarFile);
      formData.append(
        "account",
        JSON.stringify({
          name: this.createForm.name,
          phone_number: this.createForm.phone_number,
          email: this.createForm.email,
          password: this.createForm.password,
          type: this.createForm.type,
        })
      );

      this.$axios
        .post(`/accounts`, formData)
        .then((res) => {
          this.createDialog = false;
          this.createForm = {
            name: "",
            email: "",
            password: "",
            phone_number: "",
            avatarData: null,
            avatarFile: null,
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
