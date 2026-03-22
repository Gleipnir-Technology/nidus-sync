<template>
	<div class="container-fluid py-3">
		<!-- Header -->
		<div class="row mb-3">
			<slot name="header"></slot>
		</div>
		<div class="row g-3">
			<div class="col-xl-3">
				<slot name="left"></slot>
			</div>
			<div class="col-xl-6">
				<slot name="center"></slot>
			</div>
			<div class="col-xl-3">
				<slot name="right"></slot>
			</div>
		</div>
	</div>
</template>
