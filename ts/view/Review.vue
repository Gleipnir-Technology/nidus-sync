<style scoped lang="scss">
a.card-link {
	text-decoration: none;
}
</style>
<template>
	<div class="row">
		<div class="col-md-4">
			<RouterLink class="card-link" to="/_/review/pool">
				<div class="card" style="width: 100%">
					<img
						alt="overhead pool"
						class="card-img-top"
						src="/static/img/pool-overhead.jpg"
					/>
					<div class="card-body">
						<h1>Imported Pools</h1>
						<p>
							Pools that have been imported with their aerial imagery appear
							here waiting for human review before being added to the system
						</p>
					</div>
				</div>
			</RouterLink>
		</div>
		<div class="col-md-4">
			<RouterLink class="card-link" to="/_/review/site">
				<div class="card" style="width: 100%">
					<img
						alt="pool"
						class="card-img-top"
						src="/static/img/insecticide-application.jpg"
					/>
					<div class="card-body">
						<h1>Sites</h1>
						<p>
							Areas that we're tracking for potentially becoming breeding
							locations.
						</p>
					</div>
				</div>
			</RouterLink>
		</div>
	</div>
</template>

<style>
.card-link .card-img-top {
	height: 200px;
	object-fit: cover;
	object-position: center;
}

.card-link .card {
	height: 100%;
}

.card-link {
	text-decoration: none;
	color: inherit;
	display: block;
	height: 100%;
}
</style>
