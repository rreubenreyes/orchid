import Bunyan from 'bunyan';

export default {
    logger: Bunyan.createLogger({
        name: 'orchid',
        level: process.env.ORCHID_LOG_LEVEL as Bunyan.LogLevel || 'debug',
    })
}
