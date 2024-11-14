import sentry_sdk

sentry_sdk.init(
    dsn="http://1c8fca50d02d50256a8d4aeb2c16c961@192.168.56.4:9000/2",
    # Set traces_sample_rate to 1.0 to capture 100%
    # of transactions for performance monitoring.
    traces_sample_rate=1.0,
)

division_by_zero = 1 / 0
