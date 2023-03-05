package routes

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
	"github.com/joshuaschlichting/gocms/auth"
	"github.com/joshuaschlichting/gocms/config"
	"github.com/joshuaschlichting/gocms/db"
	"github.com/joshuaschlichting/gocms/middleware"
	"github.com/joshuaschlichting/gocms/presentation"
	"github.com/joshuaschlichting/gocms/templates/components"
)

func InitGetRoutes(r *chi.Mux, tmpl *template.Template, config *config.Config, queries db.Queries, middlewareMap map[string]func(http.Handler) http.Handler) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// if not logged in
		if r.Context().Value(middleware.User) == nil {
			// redirect to login
			// set Content-Type header
			w.Header().Set("Content-Type", "")
			http.Redirect(w, r, config.Auth.SignInUrl, http.StatusFound)
			return
		}

		err := tmpl.ExecuteTemplate(w, "index", map[string]interface{}{
			"sign_in_url": config.Auth.SignInUrl,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	r.Get("/getjwtandlogin", func(w http.ResponseWriter, r *http.Request) {
		// get access code from request
		code := r.URL.Query().Get("code")
		// get JWT
		accessToken, err := auth.GetAccessJWT(code)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		username, email, err := auth.GetUserInfo(accessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var tokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte(config.Auth.JWT.SecretKey), nil)
		claims := map[string]interface{}{
			"userInfo": username,
			"email":    email,
			"exp":      time.Now().Add(time.Duration(config.Auth.JWT.ExpirationTime) * time.Second).Unix(),
			"iat":      time.Now().Unix(),
			"iss":      config.Auth.JWT.Issuer,
			"aud":      config.Auth.JWT.Audience,
			"sub":      config.Auth.JWT.Subject,
			// guid for jti
			"jti": uuid.New().String(),
		}
		// jwtauth.SetExpiryIn(claims, time.Duration(config.Auth.JWT.ExpirationTime))
		jwtToken, tokenString, err := tokenAuth.Encode(claims)
		// check expiry
		if jwtToken.Expiration().Before(time.Now()) {
			log.Println("Token expired")
			http.Error(w, "Token expired", http.StatusInternalServerError)
			return
		}
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "jwt",
			Value: tokenString,
			Path:  "/",
		})
		http.Redirect(w, r, "/secure", http.StatusFound)
	})

	r.Group(func(r chi.Router) {
		jwtAuth := jwtauth.New("HS256", []byte(config.Auth.JWT.SecretKey), nil)
		r.Use(jwtauth.Verifier(jwtAuth))
		r.Use(middleware.AddClientJWTStringToCtx)
		r.Use(middleware.AuthenticateJWT)
		r.Use(middlewareMap["addUserToCtx"])

		r.Get("/edit_user_form", func(w http.ResponseWriter, r *http.Request) {
			users, err := queries.ListUsers(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var userMap []map[string]interface{}
			for _, user := range users {
				userMap = append(userMap, map[string]interface{}{
					"ID":         user.ID,
					"Name":       user.Name,
					"Email":      user.Email,
					"Attributes": string(user.Attributes),
				})

			}
			formFields := []presentation.FormField{
				{
					Name:  "ID",
					Type:  "",
					Value: "",
				},
				{
					Name:  "Name",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Email",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Attributes",
					Type:  "text",
					Value: "",
				},
			}
			p := presentation.NewPresentor(tmpl, w)
			err = p.GetEditListItemHTML("edit_user_form", "Edit User Form", "/api/user", "put", "/edit_user_form", formFields, userMap)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		r.Get("/create_user_form", func(w http.ResponseWriter, r *http.Request) {
			formFields := []presentation.FormField{
				{
					Name:  "Name",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Email",
					Type:  "text",
					Value: "",
				},
				{
					Name:  "Attributes",
					Type:  "text",
					Value: "",
				},
			}
			p := presentation.NewPresentor(tmpl, w)
			err := p.GetCreateItemFormHTML("create_user_form", "Create User Form", "/api/user", "/edit_user_form", formFields)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		r.Get("/delete_user_form", func(w http.ResponseWriter, r *http.Request) {
			users, err := queries.ListUsers(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var userMap []map[string]interface{}
			for _, user := range users {
				userMap = append(userMap, map[string]interface{}{
					"ID":         user.ID,
					"Name":       user.Name,
					"Email":      user.Email,
					"Attributes": string(user.Attributes),
				})
			}

			err = tmpl.ExecuteTemplate(w, "delete_user_form", map[string]interface{}{
				"ClickableTable": &components.ClickableTable{
					TableID:      template.JS("user_table"),
					Table:        userMap,
					CallbackFunc: template.JS("setUserInForm"),
					JavaScript: template.JS(`
						function getRowData(tableId, columnName, columnValue) {
							console.log("getRowData-> params: " + tableId  + columnName + columnValue);
							// Get the table using its ID
							const table = document.getElementById(tableId);
						
							// Get the table headers
							const headers = table.getElementsByTagName("th");
						
							// Get the index of the target column
							let targetColumnIndex;
							for (let i = 0; i < headers.length; i++) {
								if (headers[i].textContent === columnName) {
									targetColumnIndex = i;
									break;
								}
							}
						
							// Get the table rows
							const rows = table.getElementsByTagName("tr");
						
							// Loop through each row
							for (let i = 0; i < rows.length; i++) {
							const cells = rows[i].getElementsByTagName("td");
						
							// Check if the target column exists in the row
							if (cells[targetColumnIndex]) {
								// Check if the value of the target column matches the columnValue
								if (cells[targetColumnIndex].textContent === columnValue) {
								// Get the header names
								const headerNames = Array.from(headers).map(header => header.textContent);
						
								// Get the cell values
								const cellValues = Array.from(cells).map(cell => cell.textContent);
						
								// Combine the header names and cell values into an object
								const rowData = headerNames.reduce((obj, headerName, index) => {
									obj[headerName] = cellValues[index];
									return obj;
								}, {});
								return rowData;
							}}}
							// Return null if the row is not found
							return null;
						}

						function setUserInForm() {
							console.log("setUserInForm->");
							// get row data where row id == user_tableSelectedRow
							let formData = getRowData("user_table", "ID", user_tableSelectedRow);
							// get the form
							console.log(formData);
							document.getElementById("deleteUserFormID").value = formData["ID"];
							document.getElementById("deleteUserFormName").value = formData["Name"];
							document.getElementById("deleteUserFormEmail").value = formData["Email"];
							document.getElementById("deleteUserFormAttributes").value = formData["Attributes"];
							document.getElementById("deleteUserButton").setAttribute('hx-delete', '/api/user/' + formData["ID"]);
						    htmx.process(document.getElementById('deleteUserButton'));
						}`,
					),
				},
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		r.Get("/upload_form", func(w http.ResponseWriter, r *http.Request) {
			err := tmpl.ExecuteTemplate(w, "upload_form", map[string]interface{}{
				"Token":   r.Context().Value(middleware.JWTEncodedString).(string),
				"PostURL": "/upload",
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		r.Get("/securedashboard", func(w http.ResponseWriter, r *http.Request) {
			err := tmpl.ExecuteTemplate(w, "securedashboard", map[string]interface{}{})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		r.Get("/secure", func(w http.ResponseWriter, r *http.Request) {
			var username string

			if r.Context().Value(middleware.User) == nil {
				username = ""
			} else {
				username = r.Context().Value(middleware.User).(db.User).Name
			}

			err := tmpl.ExecuteTemplate(w, "index", map[string]interface{}{
				"SecureText":  username,
				"sign_in_url": config.Auth.SignInUrl,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		})

		r.Get("/sidebar", func(w http.ResponseWriter, r *http.Request) {
			err := tmpl.ExecuteTemplate(w, "sidebar", map[string]interface{}{})
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		r.Get("/messages", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			err := tmpl.ExecuteTemplate(w, "index", map[string]interface{}{
				"SecureText": "messages",
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

	})

}

var mainHtml = template.HTML(`
<div class="pagetitle">
<h1>Dashboard</h1>
<nav>
  <ol class="breadcrumb">
	<li class="breadcrumb-item"><a>Home</a></li>
	<li class="breadcrumb-item active">Dashboard</li>
  </ol>
</nav>
</div><!-- End Page Title -->
<section class="section dashboard">
<div class="row">

  <!-- Left side columns -->
  <div class="col-lg-8">
	<div class="row">

	  <!-- Sales Card -->
	  <div class="col-xxl-4 col-md-6">
		<div class="card info-card sales-card">

		  <div class="filter">
			<a class="icon" href="#" data-bs-toggle="dropdown"><i class="bi bi-three-dots"></i></a>
			<ul class="dropdown-menu dropdown-menu-end dropdown-menu-arrow">
			  <li class="dropdown-header text-start">
				<h6>Filter</h6>
			  </li>

			  <li><a class="dropdown-item" href="#">Today</a></li>
			  <li><a class="dropdown-item" href="#">This Month</a></li>
			  <li><a class="dropdown-item" href="#">This Year</a></li>
			</ul>
		  </div>

		  <div class="card-body">
			<h5 class="card-title">Sales <span>| Today</span></h5>

			<div class="d-flex align-items-center">
			  <div class="card-icon rounded-circle d-flex align-items-center justify-content-center">
				<i class="bi bi-cart"></i>
			  </div>
			  <div class="ps-3">
				<h6>145</h6>
				<span class="text-success small pt-1 fw-bold">12%</span> <span class="text-muted small pt-2 ps-1">increase</span>

			  </div>
			</div>
		  </div>

		</div>
	  </div><!-- End Sales Card -->

	  <!-- Revenue Card -->
	  <div class="col-xxl-4 col-md-6">
		<div class="card info-card revenue-card">

		  <div class="filter">
			<a class="icon" href="#" data-bs-toggle="dropdown"><i class="bi bi-three-dots"></i></a>
			<ul class="dropdown-menu dropdown-menu-end dropdown-menu-arrow">
			  <li class="dropdown-header text-start">
				<h6>Filter</h6>
			  </li>

			  <li><a class="dropdown-item" href="#">Today</a></li>
			  <li><a class="dropdown-item" href="#">This Month</a></li>
			  <li><a class="dropdown-item" href="#">This Year</a></li>
			</ul>
		  </div>

		  <div class="card-body">
			<h5 class="card-title">Revenue <span>| This Month</span></h5>

			<div class="d-flex align-items-center">
			  <div class="card-icon rounded-circle d-flex align-items-center justify-content-center">
				<i class="bi bi-currency-dollar"></i>
			  </div>
			  <div class="ps-3">
				<h6>$3,264</h6>
				<span class="text-success small pt-1 fw-bold">8%</span> <span class="text-muted small pt-2 ps-1">increase</span>

			  </div>
			</div>
		  </div>

		</div>
	  </div><!-- End Revenue Card -->

	  <!-- Customers Card -->
	  <div class="col-xxl-4 col-xl-12">

		<div class="card info-card customers-card">

		  <div class="filter">
			<a class="icon" href="#" data-bs-toggle="dropdown"><i class="bi bi-three-dots"></i></a>
			<ul class="dropdown-menu dropdown-menu-end dropdown-menu-arrow">
			  <li class="dropdown-header text-start">
				<h6>Filter</h6>
			  </li>

			  <li><a class="dropdown-item" href="#">Today</a></li>
			  <li><a class="dropdown-item" href="#">This Month</a></li>
			  <li><a class="dropdown-item" href="#">This Year</a></li>
			</ul>
		  </div>

		  <div class="card-body">
			<h5 class="card-title">Customers <span>| This Year</span></h5>

			<div class="d-flex align-items-center">
			  <div class="card-icon rounded-circle d-flex align-items-center justify-content-center">
				<i class="bi bi-people"></i>
			  </div>
			  <div class="ps-3">
				<h6>1244</h6>
				<span class="text-danger small pt-1 fw-bold">12%</span> <span class="text-muted small pt-2 ps-1">decrease</span>

			  </div>
			</div>

		  </div>
		</div>

	  </div><!-- End Customers Card -->

	  <!-- Reports -->
	  <div class="col-12">
		<div class="card">

		  <div class="filter">
			<a class="icon" href="#" data-bs-toggle="dropdown"><i class="bi bi-three-dots"></i></a>
			<ul class="dropdown-menu dropdown-menu-end dropdown-menu-arrow">
			  <li class="dropdown-header text-start">
				<h6>Filter</h6>
			  </li>

			  <li><a class="dropdown-item" href="#">Today</a></li>
			  <li><a class="dropdown-item" href="#">This Month</a></li>
			  <li><a class="dropdown-item" href="#">This Year</a></li>
			</ul>
		  </div>

		  <div class="card-body">
			<h5 class="card-title">Reports <span>/Today</span></h5>

			<!-- Line Chart -->
			<div id="reportsChart"></div>

			<script>
			  document.addEventListener("DOMContentLoaded", () => {
				new ApexCharts(document.querySelector("#reportsChart"), {
				  series: [{
					name: 'Sales',
					data: [31, 40, 28, 51, 42, 82, 56],
				  }, {
					name: 'Revenue',
					data: [11, 32, 45, 32, 34, 52, 41]
				  }, {
					name: 'Customers',
					data: [15, 11, 32, 18, 9, 24, 11]
				  }],
				  chart: {
					height: 350,
					type: 'area',
					toolbar: {
					  show: false
					},
				  },
				  markers: {
					size: 4
				  },
				  colors: ['#4154f1', '#2eca6a', '#ff771d'],
				  fill: {
					type: "gradient",
					gradient: {
					  shadeIntensity: 1,
					  opacityFrom: 0.3,
					  opacityTo: 0.4,
					  stops: [0, 90, 100]
					}
				  },
				  dataLabels: {
					enabled: false
				  },
				  stroke: {
					curve: 'smooth',
					width: 2
				  },
				  xaxis: {
					type: 'datetime',
					categories: ["2018-09-19T00:00:00.000Z", "2018-09-19T01:30:00.000Z", "2018-09-19T02:30:00.000Z", "2018-09-19T03:30:00.000Z", "2018-09-19T04:30:00.000Z", "2018-09-19T05:30:00.000Z", "2018-09-19T06:30:00.000Z"]
				  },
				  tooltip: {
					x: {
					  format: 'dd/MM/yy HH:mm'
					},
				  }
				}).render();
			  });
			</script>
			<!-- End Line Chart -->

		  </div>

		</div>
	  </div><!-- End Reports -->

	  <!-- Recent Sales -->
	  <div class="col-12">
		<div class="card recent-sales overflow-auto">

		  <div class="filter">
			<a class="icon" href="#" data-bs-toggle="dropdown"><i class="bi bi-three-dots"></i></a>
			<ul class="dropdown-menu dropdown-menu-end dropdown-menu-arrow">
			  <li class="dropdown-header text-start">
				<h6>Filter</h6>
			  </li>

			  <li><a class="dropdown-item" href="#">Today</a></li>
			  <li><a class="dropdown-item" href="#">This Month</a></li>
			  <li><a class="dropdown-item" href="#">This Year</a></li>
			</ul>
		  </div>

		  <div class="card-body">
			<h5 class="card-title">Recent Sales <span>| Today</span></h5>

			<table class="table table-borderless datatable">
			  <thead>
				<tr>
				  <th scope="col">#</th>
				  <th scope="col">Customer</th>
				  <th scope="col">Product</th>
				  <th scope="col">Price</th>
				  <th scope="col">Status</th>
				</tr>
			  </thead>
			  <tbody>
				<tr>
				  <th scope="row"><a href="#">#2457</a></th>
				  <td>Brandon Jacob</td>
				  <td><a href="#" class="text-primary">At praesentium minu</a></td>
				  <td>$64</td>
				  <td><span class="badge bg-success">Approved</span></td>
				</tr>
				<tr>
				  <th scope="row"><a href="#">#2147</a></th>
				  <td>Bridie Kessler</td>
				  <td><a href="#" class="text-primary">Blanditiis dolor omnis similique</a></td>
				  <td>$47</td>
				  <td><span class="badge bg-warning">Pending</span></td>
				</tr>
				<tr>
				  <th scope="row"><a href="#">#2049</a></th>
				  <td>Ashleigh Langosh</td>
				  <td><a href="#" class="text-primary">At recusandae consectetur</a></td>
				  <td>$147</td>
				  <td><span class="badge bg-success">Approved</span></td>
				</tr>
				<tr>
				  <th scope="row"><a href="#">#2644</a></th>
				  <td>Angus Grady</td>
				  <td><a href="#" class="text-primar">Ut voluptatem id earum et</a></td>
				  <td>$67</td>
				  <td><span class="badge bg-danger">Rejected</span></td>
				</tr>
				<tr>
				  <th scope="row"><a href="#">#2644</a></th>
				  <td>Raheem Lehner</td>
				  <td><a href="#" class="text-primary">Sunt similique distinctio</a></td>
				  <td>$165</td>
				  <td><span class="badge bg-success">Approved</span></td>
				</tr>
			  </tbody>
			</table>

		  </div>

		</div>
	  </div><!-- End Recent Sales -->

	  <!-- Top Selling -->
	  <div class="col-12">
		<div class="card top-selling overflow-auto">

		  <div class="filter">
			<a class="icon" href="#" data-bs-toggle="dropdown"><i class="bi bi-three-dots"></i></a>
			<ul class="dropdown-menu dropdown-menu-end dropdown-menu-arrow">
			  <li class="dropdown-header text-start">
				<h6>Filter</h6>
			  </li>

			  <li><a class="dropdown-item" href="#">Today</a></li>
			  <li><a class="dropdown-item" href="#">This Month</a></li>
			  <li><a class="dropdown-item" href="#">This Year</a></li>
			</ul>
		  </div>

		  <div class="card-body pb-0">
			<h5 class="card-title">Top Selling <span>| Today</span></h5>

			<table class="table table-borderless">
			  <thead>
				<tr>
				  <th scope="col">Preview</th>
				  <th scope="col">Product</th>
				  <th scope="col">Price</th>
				  <th scope="col">Sold</th>
				  <th scope="col">Revenue</th>
				</tr>
			  </thead>
			  <tbody>
				<tr>
				  <th scope="row"><a href="#"><img src="static/assets/img/product-1.jpg" alt=""></a></th>
				  <td><a href="#" class="text-primary fw-bold">Ut inventore ipsa voluptas nulla</a></td>
				  <td>$64</td>
				  <td class="fw-bold">124</td>
				  <td>$5,828</td>
				</tr>
				<tr>
				  <th scope="row"><a href="#"><img src="static/assets/img/product-2.jpg" alt=""></a></th>
				  <td><a href="#" class="text-primary fw-bold">Exercitationem similique doloremque</a></td>
				  <td>$46</td>
				  <td class="fw-bold">98</td>
				  <td>$4,508</td>
				</tr>
				<tr>
				  <th scope="row"><a href="#"><img src="static/assets/img/product-3.jpg" alt=""></a></th>
				  <td><a href="#" class="text-primary fw-bold">Doloribus nisi exercitationem</a></td>
				  <td>$59</td>
				  <td class="fw-bold">74</td>
				  <td>$4,366</td>
				</tr>
				<tr>
				  <th scope="row"><a href="#"><img src="static/assets/img/product-4.jpg" alt=""></a></th>
				  <td><a href="#" class="text-primary fw-bold">Officiis quaerat sint rerum error</a></td>
				  <td>$32</td>
				  <td class="fw-bold">63</td>
				  <td>$2,016</td>
				</tr>
				<tr>
				  <th scope="row"><a href="#"><img src="static/assets/img/product-5.jpg" alt=""></a></th>
				  <td><a href="#" class="text-primary fw-bold">Sit unde debitis delectus repellendus</a></td>
				  <td>$79</td>
				  <td class="fw-bold">41</td>
				  <td>$3,239</td>
				</tr>
			  </tbody>
			</table>

		  </div>

		</div>
	  </div><!-- End Top Selling -->

	</div>
  </div><!-- End Left side columns -->

  <!-- Right side columns -->
  <div class="col-lg-4">

	<!-- Recent Activity -->
	<div class="card">
	  <div class="filter">
		<a class="icon" href="#" data-bs-toggle="dropdown"><i class="bi bi-three-dots"></i></a>
		<ul class="dropdown-menu dropdown-menu-end dropdown-menu-arrow">
		  <li class="dropdown-header text-start">
			<h6>Filter</h6>
		  </li>

		  <li><a class="dropdown-item" href="#">Today</a></li>
		  <li><a class="dropdown-item" href="#">This Month</a></li>
		  <li><a class="dropdown-item" href="#">This Year</a></li>
		</ul>
	  </div>

	  <div class="card-body">
		<h5 class="card-title">Recent Activity <span>| Today</span></h5>

		<div class="activity">

		  <div class="activity-item d-flex">
			<div class="activite-label">32 min</div>
			<i class='bi bi-circle-fill activity-badge text-success align-self-start'></i>
			<div class="activity-content">
			  Quia quae rerum <a href="#" class="fw-bold text-dark">explicabo officiis</a> beatae
			</div>
		  </div><!-- End activity item-->

		  <div class="activity-item d-flex">
			<div class="activite-label">56 min</div>
			<i class='bi bi-circle-fill activity-badge text-danger align-self-start'></i>
			<div class="activity-content">
			  Voluptatem blanditiis blanditiis eveniet
			</div>
		  </div><!-- End activity item-->

		  <div class="activity-item d-flex">
			<div class="activite-label">2 hrs</div>
			<i class='bi bi-circle-fill activity-badge text-primary align-self-start'></i>
			<div class="activity-content">
			  Voluptates corrupti molestias voluptatem
			</div>
		  </div><!-- End activity item-->

		  <div class="activity-item d-flex">
			<div class="activite-label">1 day</div>
			<i class='bi bi-circle-fill activity-badge text-info align-self-start'></i>
			<div class="activity-content">
			  Tempore autem saepe <a href="#" class="fw-bold text-dark">occaecati voluptatem</a> tempore
			</div>
		  </div><!-- End activity item-->

		  <div class="activity-item d-flex">
			<div class="activite-label">2 days</div>
			<i class='bi bi-circle-fill activity-badge text-warning align-self-start'></i>
			<div class="activity-content">
			  Est sit eum reiciendis exercitationem
			</div>
		  </div><!-- End activity item-->

		  <div class="activity-item d-flex">
			<div class="activite-label">4 weeks</div>
			<i class='bi bi-circle-fill activity-badge text-muted align-self-start'></i>
			<div class="activity-content">
			  Dicta dolorem harum nulla eius. Ut quidem quidem sit quas
			</div>
		  </div><!-- End activity item-->

		</div>

	  </div>
	</div><!-- End Recent Activity -->

	<!-- Budget Report -->
	<div class="card">
	  <div class="filter">
		<a class="icon" href="#" data-bs-toggle="dropdown"><i class="bi bi-three-dots"></i></a>
		<ul class="dropdown-menu dropdown-menu-end dropdown-menu-arrow">
		  <li class="dropdown-header text-start">
			<h6>Filter</h6>
		  </li>

		  <li><a class="dropdown-item" href="#">Today</a></li>
		  <li><a class="dropdown-item" href="#">This Month</a></li>
		  <li><a class="dropdown-item" href="#">This Year</a></li>
		</ul>
	  </div>

	  <div class="card-body pb-0">
		<h5 class="card-title">Budget Report <span>| This Month</span></h5>

		<div id="budgetChart" style="min-height: 400px;" class="echart"></div>

		<script>
		  document.addEventListener("DOMContentLoaded", () => {
			var budgetChart = echarts.init(document.querySelector("#budgetChart")).setOption({
			  legend: {
				data: ['Allocated Budget', 'Actual Spending']
			  },
			  radar: {
				// shape: 'circle',
				indicator: [{
					name: 'Sales',
					max: 6500
				  },
				  {
					name: 'Administration',
					max: 16000
				  },
				  {
					name: 'Information Technology',
					max: 30000
				  },
				  {
					name: 'Customer Support',
					max: 38000
				  },
				  {
					name: 'Development',
					max: 52000
				  },
				  {
					name: 'Marketing',
					max: 25000
				  }
				]
			  },
			  series: [{
				name: 'Budget vs spending',
				type: 'radar',
				data: [{
					value: [4200, 3000, 20000, 35000, 50000, 18000],
					name: 'Allocated Budget'
				  },
				  {
					value: [5000, 14000, 28000, 26000, 42000, 21000],
					name: 'Actual Spending'
				  }
				]
			  }]
			});
		  });
		</script>

	  </div>
	</div><!-- End Budget Report -->

	<!-- Website Traffic -->
	<div class="card">
	  <div class="filter">
		<a class="icon" href="#" data-bs-toggle="dropdown"><i class="bi bi-three-dots"></i></a>
		<ul class="dropdown-menu dropdown-menu-end dropdown-menu-arrow">
		  <li class="dropdown-header text-start">
			<h6>Filter</h6>
		  </li>

		  <li><a class="dropdown-item" href="#">Today</a></li>
		  <li><a class="dropdown-item" href="#">This Month</a></li>
		  <li><a class="dropdown-item" href="#">This Year</a></li>
		</ul>
	  </div>

	  <div class="card-body pb-0">
		<h5 class="card-title">Website Traffic <span>| Today</span></h5>

		<div id="trafficChart" style="min-height: 400px;" class="echart"></div>

		<script>
		  document.addEventListener("DOMContentLoaded", () => {
			echarts.init(document.querySelector("#trafficChart")).setOption({
			  tooltip: {
				trigger: 'item'
			  },
			  legend: {
				top: '5%',
				left: 'center'
			  },
			  series: [{
				name: 'Access From',
				type: 'pie',
				radius: ['40%', '70%'],
				avoidLabelOverlap: false,
				label: {
				  show: false,
				  position: 'center'
				},
				emphasis: {
				  label: {
					show: true,
					fontSize: '18',
					fontWeight: 'bold'
				  }
				},
				labelLine: {
				  show: false
				},
				data: [{
					value: 1048,
					name: 'Search Engine'
				  },
				  {
					value: 735,
					name: 'Direct'
				  },
				  {
					value: 580,
					name: 'Email'
				  },
				  {
					value: 484,
					name: 'Union Ads'
				  },
				  {
					value: 300,
					name: 'Video Ads'
				  }
				]
			  }]
			});
		  });
		</script>

	  </div>
	</div><!-- End Website Traffic -->

	<!-- News & Updates Traffic -->
	<div class="card">
	  <div class="filter">
		<a class="icon" href="#" data-bs-toggle="dropdown"><i class="bi bi-three-dots"></i></a>
		<ul class="dropdown-menu dropdown-menu-end dropdown-menu-arrow">
		  <li class="dropdown-header text-start">
			<h6>Filter</h6>
		  </li>

		  <li><a class="dropdown-item" href="#">Today</a></li>
		  <li><a class="dropdown-item" href="#">This Month</a></li>
		  <li><a class="dropdown-item" href="#">This Year</a></li>
		</ul>
	  </div>

	  <div class="card-body pb-0">
		<h5 class="card-title">News &amp; Updates <span>| Today</span></h5>

		<div class="news">
		  <div class="post-item clearfix">
			<img src="static/assets/img/news-1.jpg" alt="">
			<h4><a href="#">Nihil blanditiis at in nihil autem</a></h4>
			<p>Sit recusandae non aspernatur laboriosam. Quia enim eligendi sed ut harum...</p>
		  </div>

		  <div class="post-item clearfix">
			<img src="static/assets/img/news-2.jpg" alt="">
			<h4><a href="#">Quidem autem et impedit</a></h4>
			<p>Illo nemo neque maiores vitae officiis cum eum turos elan dries werona nande...</p>
		  </div>

		  <div class="post-item clearfix">
			<img src="static/assets/img/news-3.jpg" alt="">
			<h4><a href="#">Id quia et et ut maxime similique occaecati ut</a></h4>
			<p>Fugiat voluptas vero eaque accusantium eos. Consequuntur sed ipsam et totam...</p>
		  </div>

		  <div class="post-item clearfix">
			<img src="static/assets/img/news-4.jpg" alt="">
			<h4><a href="#">Laborum corporis quo dara net para</a></h4>
			<p>Qui enim quia optio. Eligendi aut asperiores enim repellendusvel rerum cuderouter...</p>
		  </div>

		  <div class="post-item clearfix">
			<img src="static/assets/img/news-5.jpg" alt="">
			<h4><a href="#">Et dolores corrupti quae illo quod dolor</a></h4>
			<p>Odit ut eveniet modi reiciendis. Atque cupiditate libero beatae dignissimos eius...</p>
		  </div>

		</div><!-- End sidebar recent posts-->

	  </div>
	</div><!-- End News & Updates -->

  </div><!-- End Right side columns -->

</div>
</section>
`)
