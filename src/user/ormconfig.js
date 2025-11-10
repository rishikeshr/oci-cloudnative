/**
 * TypeOrm configuration file
 * By using this file, we are able to incorporate CLI usage, and
 * the schema synchronization into CI pipelines
 * @see https://typeorm.io/#/using-ormconfig
 */

// read env
const {
  POSTGRES_HOST,
  POSTGRES_PORT,
  POSTGRES_DB,
  POSTGRES_USER,
  POSTGRES_PASSWORD,
  NODE_ENV,
} = process.env;

// determine opts
const prod = /^prod/i.test(NODE_ENV || '');
const useExt = prod ? 'js' : 'ts';

module.exports = {
  type: 'postgres',
  host: POSTGRES_HOST || 'localhost',
  port: parseInt(POSTGRES_PORT || '5432', 10),
  database: POSTGRES_DB || 'mushop_user',
  username: POSTGRES_USER || 'mushop',
  password: POSTGRES_PASSWORD || 'mushop',
  entities: [
    // use glob and prevent dupes in development when both `src` and `dist` exist
    `**/*.entity.${useExt}`,
  ],
  synchronize: false,
  logging: !prod,
};
